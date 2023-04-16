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

	SaveBinaryData(ctx context.Context, name string, data []byte, userID string) error
	GetBinaryData(ctx context.Context, name, userID string) ([]byte, error)

	SaveCardData(ctx context.Context, name string, cardData entity.CardDataDTO, userID string) error
	GetCardData(ctx context.Context, name, userID string) (entity.CardDataDTO, error)
}
