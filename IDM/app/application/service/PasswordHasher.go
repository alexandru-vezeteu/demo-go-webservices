package service

// PasswordHasher defines the interface for password hashing operations
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
}
