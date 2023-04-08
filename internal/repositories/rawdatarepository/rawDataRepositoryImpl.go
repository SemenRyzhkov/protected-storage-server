package rawdatarepository

import (
	"context"
	"database/sql"

	"github.com/omeid/pgerror"
	"github.com/rs/zerolog/log"

	"protected-storage-server/internal/entity/myerrors"
)

const (
	insertRawDataQuery = "" +
		"INSERT INTO public.raw_data (name, data, user_id) " +
		"VALUES ($1, $2, $3)"
	getRawDataQuery = "" +
		"SELECT data FROM public.raw_data " +
		"WHERE user_id=$1 AND name=$2"
)

type rawDataRepositoryImpl struct {
	db *sql.DB
}

// New конструктор UserRepository
func New(db *sql.DB) RawDataRepository {
	return &rawDataRepositoryImpl{
		db: db,
	}
}

// Save сохранение произвольных текстовых данных
func (r *rawDataRepositoryImpl) Save(ctx context.Context, userID, name string, data []byte) error {
	log.Info().Msgf("rawdatarepository: save raw data with name %s for user with ID %s to db", name, userID)
	_, err := r.db.ExecContext(ctx, insertRawDataQuery, name, data, userID)

	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewDataViolationError(name, err)
		}
		return err
	}
	return nil
}

// GetByName получение произвольных текстовых данных
func (r *rawDataRepositoryImpl) GetByName(ctx context.Context, userID, name string) ([]byte, error) {
	var data []byte
	log.Info().Msgf("rawdatarepository: get raw data with name %s for user with ID %s to db", name, userID)
	row := r.db.QueryRowContext(ctx, getRawDataQuery, userID, name)
	err := row.Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, myerrors.NewNotFoundError(name, err)
		}
		return nil, err
	}
	return data, nil
}
