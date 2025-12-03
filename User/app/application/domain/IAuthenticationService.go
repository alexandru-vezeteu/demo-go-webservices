package domain

// Credentials represents user credentials for authentication
type Credentials struct {
	Email    string
	Password string
}

// AuthToken represents an authentication token
type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // Expiration time in seconds
	TokenType    string
}

// AuthenticatedUser represents an authenticated user with token information
type AuthenticatedUser struct {
	User  *User
	Token *AuthToken
}

// IAuthenticationService defines the contract for authentication operations
type IAuthenticationService interface {
	// Authenticate verifies user credentials and returns an authenticated user with tokens
	Authenticate(credentials *Credentials) (*AuthenticatedUser, error)

	// ValidateToken validates an access token and returns the associated user
	ValidateToken(token string) (*User, error)

	// RefreshToken generates a new access token using a refresh token
	RefreshToken(refreshToken string) (*AuthToken, error)

	// RevokeToken invalidates a token (logout)
	RevokeToken(token string) error

	// ChangePassword changes the password for a user
	ChangePassword(userID int, oldPassword, newPassword string) error
}
