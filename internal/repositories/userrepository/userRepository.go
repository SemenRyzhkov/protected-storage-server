package userrepository

import (
	"context"
)

type UserRepository interface {
	Save(ctx context.Context, userID, login, password string) error
}
