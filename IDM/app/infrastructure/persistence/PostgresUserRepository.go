package persistence

import (
	"errors"
	"idmService/domain"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *PostgresUserRepository) Update(user *domain.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *PostgresUserRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.User{}, id)
	return result.Error
}

func (r *PostgresUserRepository) MigrateSchema() error {
	return r.db.AutoMigrate(&domain.User{})
}
