package rawdatarepository

import (
	"context"
)

type RawDataRepository interface {
	Save(ctx context.Context, userID, name string, data []byte) error
	GetByName(ctx context.Context, userID, name string) ([]byte, error)
}
