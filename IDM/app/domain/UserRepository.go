package domain

import "context"

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)

	FindByID(ctx context.Context, id uint) (*User, error)

	Create(ctx context.Context, user *User) error

	Update(ctx context.Context, user *User) error

	Delete(ctx context.Context, id uint) error

	MigrateSchema() error
}
