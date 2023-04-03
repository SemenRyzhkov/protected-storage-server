package userrepository

import (
	"context"
	"database/sql"

	"github.com/omeid/pgerror"
	"github.com/rs/zerolog/log"

	"protected-storage-server/internal/entity"
	"protected-storage-server/internal/entity/myerrors"
)

const (
	insertUserQuery = "" +
		"INSERT INTO public.users (id, login, password) " +
		"VALUES ($1, $2, $3)"
	findUserByLoginQuery = "" +
		"SELECT id, login, password FROM public.users " +
		"WHERE login=$1"
)

type userRepositoryImpl struct {
	db *sql.DB
}

// New конструктор UserRepository
func New(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Save сохранение юзера
func (r *userRepositoryImpl) Save(ctx context.Context, userID, login, password string) error {
	log.Info().Msgf("userrepository: save user with ID %s and login %s to db", userID, login)
	_, err := r.db.ExecContext(ctx, insertUserQuery, userID, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			log.Error().Msgf("userrepository: user with login %s already exists")
			return myerrors.NewUserViolationError(login, err)
		}
	}
	return nil
}

// FindByLogin поиск юзера в базе по логину
func (r *userRepositoryImpl) FindByLogin(ctx context.Context, login string) (entity.UserDTO, error) {
	var user entity.UserDTO
	log.Info().Msgf("userrepository: find user with login %s in db", login)
	row := r.db.QueryRowContext(ctx, findUserByLoginQuery, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.UserDTO{}, err
	}
	return user, nil
}
