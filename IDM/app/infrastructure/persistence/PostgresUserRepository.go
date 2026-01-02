package persistence

import (
	"context"
	"errors"

	"idmService/application/domain"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var gormUser GormUser
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &domain.InternalError{Operation: "find user by email", Err: result.Error}
	}

	return gormUser.ToDomain(), nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var gormUser GormUser
	result := r.db.WithContext(ctx).First(&gormUser, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, &domain.InternalError{Operation: "find user by ID", Err: result.Error}
	}

	return gormUser.ToDomain(), nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	gormUser := FromDomainUser(user)
	result := r.db.WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		return &domain.InternalError{Operation: "create user", Err: result.Error}
	}
	user.ID = gormUser.ID
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	gormUser := FromDomainUser(user)
	result := r.db.WithContext(ctx).Save(gormUser)
	if result.Error != nil {
		return &domain.InternalError{Operation: "update user", Err: result.Error}
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&GormUser{}, id)
	if result.Error != nil {
		return &domain.InternalError{Operation: "delete user", Err: result.Error}
	}
	return nil
}

func (r *PostgresUserRepository) MigrateSchema() error {
	if err := r.db.AutoMigrate(&GormUser{}); err != nil {
		return &domain.InternalError{Operation: "migrate schema", Err: err}
	}
	return nil
}

func (r *PostgresUserRepository) SeedServiceAccount(email, password string) error {
	ctx := context.Background()

	existing, err := r.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	serviceAccount := &domain.User{
		Email:  email,
		Parola: password,
		Rol:    domain.RoleServiceClient,
	}

	return r.Create(ctx, serviceAccount)
}
