package grpcserver

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protected-storage-server/internal/entity/myerrors"
	"protected-storage-server/internal/security"
	"protected-storage-server/internal/service/userservice"
	"protected-storage-server/proto"
)

// Server сервер
type Server struct {
	proto.UnimplementedGrpcServiceServer
	userService userservice.UserService
	jwtManager  *security.JWTManager
}

// NewServer конструктор.
func NewServer(userService userservice.UserService, jwtHelper *security.JWTManager) *Server {
	return &Server{userService: userService, jwtManager: jwtHelper}
}

// CreateUser выполняет сохранение нового пользователя, генерит токен и отдает в теле респонса
func (s *Server) CreateUser(ctx context.Context, in *proto.UserRegisterRequest) (*proto.AuthorizedResponse, error) {
	login := in.Login
	password := in.Password
	userID := uuid.New().String()

	err := s.userService.Create(ctx, login, password, userID)
	if err != nil {
		var uv *myerrors.UserViolationError
		if errors.As(err, &uv) {
			return nil, status.Errorf(codes.Unauthenticated, uv.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	token, err := s.jwtManager.GenerateJWT(userID, login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.AuthorizedResponse{Token: token}, nil
}

// LoginUser выполняет авторизацию существующего пользователя, генерит токен и отдает в теле респонса
func (s *Server) LoginUser(ctx context.Context, in *proto.UserAuthorizedRequest) (*proto.AuthorizedResponse, error) {
	login := in.Login
	password := in.Password

	userID, err := s.userService.Login(ctx, login, password)
	if err != nil {
		var ip *myerrors.InvalidPasswordError
		if errors.As(err, &ip) {
			return nil, status.Errorf(codes.Unauthenticated, ip.Error())
		}
		return nil, status.Errorf(codes.Internal, "user with login %s not found", login)
	}

	token, err := s.jwtManager.GenerateJWT(userID, login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.AuthorizedResponse{Token: token}, nil
}

// SaveRawData выполняет сохранение текстовой информации для авторизованного пользователя
func (s *Server) SaveRawData(ctx context.Context, in *proto.SaveRawDataRequest) (*proto.ErrorResponse, error) {

	return nil, nil
}
