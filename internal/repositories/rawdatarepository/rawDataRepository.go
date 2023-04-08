package rawdatarepository

import (
	"context"
)

type RawDataRepository interface {
	Save(ctx context.Context, userID, name, data string) error
	GetByName(ctx context.Context, userID, name string) (string, error)
}
