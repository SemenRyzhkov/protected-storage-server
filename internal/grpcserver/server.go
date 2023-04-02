package grpcserver

import (
	"protected-storage-server/proto"
)

// Server сервер
type Server struct {
	proto.UnimplementedGrpcServiceServer
}

// NewServer конструктор.
func NewServer() *Server {
	return &Server{}
}
