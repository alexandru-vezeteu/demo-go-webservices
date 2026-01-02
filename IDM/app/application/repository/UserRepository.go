package repository

import (
	"context"

	"idmService/application/domain"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	FindByID(ctx context.Context, id uint) (*domain.User, error)

	Create(ctx context.Context, user *domain.User) error

	Update(ctx context.Context, user *domain.User) error

	Delete(ctx context.Context, id uint) error

	MigrateSchema() error
}
