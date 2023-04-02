package grpcserver

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"protected-storage-server/internal/entity/myerrors"
	"protected-storage-server/internal/security"
	"protected-storage-server/internal/service/userservice"
	"protected-storage-server/proto"
)

// Server сервер
type Server struct {
	proto.UnimplementedGrpcServiceServer
	userService userservice.UserService
	jwtHelper   *security.JwtHelper
}

// NewServer конструктор.
func NewServer(userService userservice.UserService, jwtHelper *security.JwtHelper) *Server {
	return &Server{userService: userService, jwtHelper: jwtHelper}
}

// CreateUser выполняет сохранение нового пользователя, генерацию и добавление токена в метаданные
func (s *Server) CreateUser(ctx context.Context, in *proto.UserAuthorizeRequest) (*proto.ErrorResponse, error) {
	login := in.Login
	password := in.Password
	userID := uuid.New().String()

	err := s.userService.Create(ctx, login, password, userID)
	if err != nil {
		var uv *myerrors.UserViolationError
		if errors.As(err, &uv) {
			return &proto.ErrorResponse{Error: uv.Error()}, nil
		}
		return &proto.ErrorResponse{Error: "Internal Server Error"}, nil
	}

	token, err := s.jwtHelper.GenerateJWT(userID)
	if err != nil {
		return &proto.ErrorResponse{Error: "Internal Server Error"}, nil
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "token", token)

	return nil, nil
}
