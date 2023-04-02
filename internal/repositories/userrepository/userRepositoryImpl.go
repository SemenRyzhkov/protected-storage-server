package userrepository

import (
	"context"
	"database/sql"

	"github.com/omeid/pgerror"

	"protected-storage-server/internal/entity/myerrors"
)

const (
	insertUserQuery = "" +
		"INSERT INTO public.users (id, login, password) " +
		"VALUES ($1, $2, $3)"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Save(ctx context.Context, userID, login, password string) error {
	_, err := r.db.ExecContext(ctx, insertUserQuery, userID, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewUserViolationError(login, err)
		}
	}
	return nil
}
