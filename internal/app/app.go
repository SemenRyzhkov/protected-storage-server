package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"protected-storage-server/internal/config"
	"protected-storage-server/internal/grpcserver"
	"protected-storage-server/proto"
)

// GRPCApp запускает GRPC приложение.
type GRPCApp struct {
	GRPCServer *grpc.Server
}

// NewGRPC конструктор GRPCApp
func NewGRPC(cfg config.Config) (*GRPCApp, error) {
	log.Println("creating router")

	serverImpl := grpcserver.NewServer()

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
	fmt.Println("Сервер GRPc начал работу")

	// получаем запрос gRPC
	if err := app.GRPCServer.Serve(listen); err != nil {
		return err
	}
	return nil
}
