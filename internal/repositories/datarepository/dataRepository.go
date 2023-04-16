package datarepository

import (
	"context"

	"protected-storage-server/internal/entity"
)

type RawDataRepository interface {
	Save(ctx context.Context, userID, name string, data []byte, dataType entity.DataType) error
	GetByNameAndTypeAndUserID(ctx context.Context, userID, name string, dataType entity.DataType) ([]byte, error)
}
