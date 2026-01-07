package repository

import (
	"context"
	"userService/application/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.User, error)
	Delete(ctx context.Context, id int) (*domain.User, error)

	GetUsersByEventID(ctx context.Context, eventID int) ([]*domain.User, error)
	GetUsersByPacketID(ctx context.Context, packetID int) ([]*domain.User, error)
}
