package app

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"protected-storage-server/internal/config"
	"protected-storage-server/internal/grpcserver"
	"protected-storage-server/internal/repositories"
	"protected-storage-server/internal/repositories/userrepository"
	"protected-storage-server/internal/security"
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
	userService := userservice.New(userRepository)
	jwtHelper, err := security.New(cfg.Key)
	if err != nil {
		return nil, err
	}

	serverImpl := grpcserver.NewServer(userService, jwtHelper)

	s := grpc.NewServer()

	proto.RegisterGrpcServiceServer(s, serverImpl)

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
