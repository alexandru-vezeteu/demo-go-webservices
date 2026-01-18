package service

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
}
