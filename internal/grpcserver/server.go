package grpcserver

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"protected-storage-server/internal/entity"
	"protected-storage-server/internal/entity/myerrors"
	"protected-storage-server/internal/security"
	"protected-storage-server/internal/service/dataservice"
	"protected-storage-server/internal/service/userservice"
	"protected-storage-server/proto"
)

// Server сервер
type Server struct {
	proto.UnimplementedGrpcServiceServer
	userService    userservice.UserService
	storageService dataservice.StorageService
	jwtManager     *security.JWTManager
}

// NewServer конструктор.
func NewServer(userService userservice.UserService, storageService dataservice.StorageService, jwtHelper *security.JWTManager) *Server {
	return &Server{userService: userService, jwtManager: jwtHelper, storageService: storageService}
}

// CreateUser эндпойнт сохранения нового пользователя, генерит токен и отдает в теле респонса
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

// LoginUser эндпойнт авторизации существующего пользователя, генерит токен и отдает в теле респонса
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

// SaveRawData эндпойнт сохранения произвольной текстовой информации для авторизованного пользователя
func (s *Server) SaveRawData(ctx context.Context, in *proto.SaveRawDataRequest) (*proto.ErrorResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	err = s.storageService.SaveRawData(ctx, in.Name, in.Data, userID)
	if err != nil {
		var dv *myerrors.DataViolationError
		if errors.As(err, &dv) {
			return nil, status.Errorf(codes.AlreadyExists, dv.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.ErrorResponse{}, nil
}

// GetRawData эндпойнт получения текстовой информации по названию для авторизованного пользователя
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

// SaveLoginWithPassword эндпойнт сохранения логина и пароля для авторизованного пользователя
func (s *Server) SaveLoginWithPassword(ctx context.Context, in *proto.SaveLoginWithPasswordRequest) (*proto.ErrorResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	err = s.storageService.SaveLoginWithPassword(ctx, in.Name, in.Login, in.Password, userID)
	if err != nil {
		var dv *myerrors.DataViolationError
		if errors.As(err, &dv) {
			return nil, status.Errorf(codes.AlreadyExists, dv.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.ErrorResponse{}, nil
}

// GetLoginWithPassword эндпойнт получения логина и пароля по названию для авторизованного пользователя
func (s *Server) GetLoginWithPassword(ctx context.Context, in *proto.GetLoginWithPasswordRequest) (*proto.GetLoginWithPasswordResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	data, err := s.storageService.GetLoginWithPassword(ctx, in.Name, userID)
	if err != nil {
		var nf *myerrors.NotFoundError
		if errors.As(err, &nf) {
			return nil, status.Errorf(codes.NotFound, nf.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.GetLoginWithPasswordResponse{Login: data.Login, Password: data.Password}, nil
}

// SaveBinaryData эндпойнт сохранения произвольных бинарных данных для авторизованного пользователя
func (s *Server) SaveBinaryData(ctx context.Context, in *proto.SaveBinaryDataRequest) (*proto.ErrorResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	err = s.storageService.SaveBinaryData(ctx, in.Name, in.Data, userID)
	if err != nil {
		var dv *myerrors.DataViolationError
		if errors.As(err, &dv) {
			return nil, status.Errorf(codes.AlreadyExists, dv.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.ErrorResponse{}, nil
}

// GetBinaryData эндпойнт получения произвольных бинарных данных по названию для авторизованного пользователя
func (s *Server) GetBinaryData(ctx context.Context, in *proto.GetBinaryDataRequest) (*proto.GetBinaryDataResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	data, err := s.storageService.GetBinaryData(ctx, in.Name, userID)
	if err != nil {
		var nf *myerrors.NotFoundError
		if errors.As(err, &nf) {
			return nil, status.Errorf(codes.NotFound, nf.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.GetBinaryDataResponse{Data: data}, nil
}

// SaveCardData эндпойнт сохранения данных банковской карты для авторизованного пользователя
func (s *Server) SaveCardData(ctx context.Context, in *proto.SaveCardDataRequest) (*proto.ErrorResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	card := entity.CardDataDTO{
		Number:     in.Number,
		Month:      in.Month,
		Year:       in.Year,
		CardHolder: in.CardHolder,
	}

	err = s.storageService.SaveCardData(ctx, in.Name, card, userID)
	if err != nil {
		var dv *myerrors.DataViolationError
		if errors.As(err, &dv) {
			return nil, status.Errorf(codes.AlreadyExists, dv.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.ErrorResponse{}, nil
}

// GetCardData эндпойнт получения данных банковской карты по названию для авторизованного пользователя
func (s *Server) GetCardData(ctx context.Context, in *proto.GetCardDataRequest) (*proto.GetCardDataResponse, error) {
	userID, err := s.jwtManager.ExtractUserID(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	data, err := s.storageService.GetCardData(ctx, in.Name, userID)
	if err != nil {
		var nf *myerrors.NotFoundError
		if errors.As(err, &nf) {
			return nil, status.Errorf(codes.NotFound, nf.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.GetCardDataResponse{Number: data.Number, Month: data.Number, Year: data.Year, CardHolder: data.CardHolder}, nil
}
