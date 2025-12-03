package domain

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleOwnerEvent UserRole = "owner-event"
	RoleClient     UserRole = "client"
)

type User struct {
	ID     uint     `gorm:"primaryKey;autoIncrement"`
	Email  string   `gorm:"type:varchar(255);uniqueIndex;not null"`
	Parola string   `gorm:"type:varchar(255);not null"`
	Rol    UserRole `gorm:"type:varchar(50);not null"`
}

func (User) TableName() string {
	return "users"
}

func ValidateRole(role string) bool {
	switch UserRole(role) {
	case RoleAdmin, RoleOwnerEvent, RoleClient:
		return true
	default:
		return false
	}
}
