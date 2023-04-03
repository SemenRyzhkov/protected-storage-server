package userservice

import (
	"context"
	"encoding/base64"

	"github.com/rs/zerolog/log"

	"protected-storage-server/internal/entity/myerrors"
	"protected-storage-server/internal/repositories/userrepository"
)

var _ UserService = &userServiceImpl{}

type userServiceImpl struct {
	userRepository userrepository.UserRepository
}

// Create метод для регистрации пользователя
func (u userServiceImpl) Create(ctx context.Context, login, password, userID string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	log.Info().Msgf("save user with ID %s and login %s to db", userID, login)
	return u.userRepository.Save(ctx, userID, login, encodedPassword)
}

// Login метод для авторизации
func (u userServiceImpl) Login(ctx context.Context, login, password string) (string, error) {
	log.Info().Msgf("userservice: find user with login %s in db", login)
	foundUser, err := u.userRepository.FindByLogin(ctx, login)
	if err != nil {
		log.Error().Msgf("user with login %s not found", login)
		return "", err
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(foundUser.Password)
	if err != nil {
		return "", err
	}

	if password != string(decodedPassword) {
		log.Error().Msgf("userservice: password %s is invalid", password)
		return "", &myerrors.InvalidPasswordError{Password: password}
	}

	return foundUser.ID, nil
}

// New конструктор UserService
func New(userRepository userrepository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
