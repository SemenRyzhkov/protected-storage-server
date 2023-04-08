package storageservice

import (
	"context"
)

type StorageService interface {
	SaveRawData(ctx context.Context, name, data, userID string) error
	GetRawData(ctx context.Context, name, userID string) (string, error)
}
