package grpcserver

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protected-storage-server/internal/entity/myerrors"
	"protected-storage-server/internal/security"
	"protected-storage-server/internal/service/storageservice"
	"protected-storage-server/internal/service/userservice"
	"protected-storage-server/proto"
)

// Server сервер
type Server struct {
	proto.UnimplementedGrpcServiceServer
	userService    userservice.UserService
	storageService storageservice.StorageService
	jwtManager     *security.JWTManager
}

// NewServer конструктор.
func NewServer(userService userservice.UserService, storageService storageservice.StorageService, jwtHelper *security.JWTManager) *Server {
	return &Server{userService: userService, jwtManager: jwtHelper, storageService: storageService}
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
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	err = s.storageService.SaveRawData(ctx, in.Name, in.Data, userID)
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}

	return &proto.ErrorResponse{}, nil
}

// GetRawData выполняет получение текстовой информации по названию для авторизованного пользователя
func (s *Server) GetRawData(ctx context.Context, in *proto.GetRawDataRequest) (*proto.GetRawDataResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	data, err := s.storageService.GetRawData(ctx, in.Name, userID)
	if err != nil {
		var nf *myerrors.NotFoundError
		if errors.As(err, &nf) {
			return nil, status.Errorf(codes.NotFound, nf.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.GetRawDataResponse{Data: data}, nil
}
