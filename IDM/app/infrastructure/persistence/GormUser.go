package persistence

import "idmService/application/domain"

type GormUser struct {
	ID     uint            `gorm:"primaryKey;autoIncrement"`
	Email  string          `gorm:"type:varchar(255);uniqueIndex;not null"`
	Parola string          `gorm:"type:varchar(255);not null"`
	Rol    domain.UserRole `gorm:"type:varchar(50);not null"`
}

func (GormUser) TableName() string {
	return "users"
}

func (g *GormUser) ToDomain() *domain.User {
	return &domain.User{
		ID:     g.ID,
		Email:  g.Email,
		Parola: g.Parola,
		Rol:    g.Rol,
	}
}

func FromDomainUser(u *domain.User) *GormUser {
	return &GormUser{
		ID:     u.ID,
		Email:  u.Email,
		Parola: u.Parola,
		Rol:    u.Rol,
	}
}
