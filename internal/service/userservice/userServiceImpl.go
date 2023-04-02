package userservice

import (
	"context"
	"encoding/base64"

	"protected-storage-server/internal/repositories/userrepository"
	"protected-storage-server/internal/security"
)

var _ UserService = &userServiceImpl{}

type userServiceImpl struct {
	userRepository userrepository.UserRepository
}

func (u userServiceImpl) Create(ctx context.Context, login, password, userID string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

	return u.userRepository.Save(ctx, userID, login, encodedPassword)
}

func New(userRepository userrepository.UserRepository, jwtHelper *security.JwtHelper) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
