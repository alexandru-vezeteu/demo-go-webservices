package domain

type UserRole string

const (
	RoleAdmin         UserRole = "admin"
	RoleOwnerEvent    UserRole = "owner-event"
	RoleClient        UserRole = "client"
	RoleServiceClient UserRole = "serviciu_clienti"
)

type User struct {
	ID     uint
	Email  string
	Parola string
	Rol    UserRole
}

func ValidateRole(role string) bool {
	switch UserRole(role) {
	case RoleAdmin, RoleOwnerEvent, RoleClient, RoleServiceClient:
		return true
	default:
		return false
	}
}
