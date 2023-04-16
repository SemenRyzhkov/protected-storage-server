package dataservice

import (
	"context"

	"protected-storage-server/internal/entity"
)

type StorageService interface {
	SaveRawData(ctx context.Context, name, data, userID string) error
	GetRawData(ctx context.Context, name, userID string) (string, error)

	SaveLoginWithPassword(ctx context.Context, name, login, password, userID string) error
	GetLoginWithPassword(ctx context.Context, name, userID string) (entity.CredentialsDTO, error)
}
