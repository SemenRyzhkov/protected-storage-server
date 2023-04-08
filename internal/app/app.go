package app

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"protected-storage-server/internal/app/interceptors"
	"protected-storage-server/internal/config"
	"protected-storage-server/internal/grpcserver"
	"protected-storage-server/internal/repositories"
	"protected-storage-server/internal/repositories/rawdatarepository"
	"protected-storage-server/internal/repositories/userrepository"
	"protected-storage-server/internal/security"
	"protected-storage-server/internal/service/storageservice"
	"protected-storage-server/internal/service/userservice"
	"protected-storage-server/proto"
)

// GRPCApp запускает GRPC приложение.
type GRPCApp struct {
	GRPCServer *grpc.Server
}

// NewGRPC конструктор GRPCApp
func NewGRPC(cfg config.Config) (*GRPCApp, error) {
	log.Println("creating server")

	db, err := repositories.InitDB(cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}

	userRepository := userrepository.New(db)
	rawDataRepository := rawdatarepository.New(db)
	userService := userservice.New(userRepository)
	jwtManager, err := security.NewJWTManager(cfg.Key, cfg.TokenDuration)
	storageService := storageservice.New(rawDataRepository)
	if err != nil {
		return nil, err
	}

	serverImpl := grpcserver.NewServer(userService, storageService, jwtManager)

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)
	logInterceptor := interceptors.NewLogInterceptor()

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logInterceptor.Unary(),
			authInterceptor.Unary(),
		),
	)

	proto.RegisterGrpcServiceServer(s, serverImpl)
	reflection.Register(s)

	return &GRPCApp{
		GRPCServer: s,
	}, nil

}

// Run запуск сервера
func (app *GRPCApp) Run(cfg config.Config) error {
	listen, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		return err
	}
	log.Println("Start GRPc-server")

	// получаем запрос gRPC
	if err := app.GRPCServer.Serve(listen); err != nil {
		return err
	}
	return nil
}
