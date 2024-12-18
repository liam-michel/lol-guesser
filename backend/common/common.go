package common

import "fmt"

var (
	ErrTokenExpired         = fmt.Errorf("token has expired")
	ErrInvalidSignature     = fmt.Errorf("invalid token signature")
	ErrInvalidSigningMethod = fmt.Errorf("invalid signing method: expected HS256")
	ErrMalformedToken       = fmt.Errorf("malformed token")
	ErrInvalidClaims        = fmt.Errorf("invalid token claims")
)

type Username struct {
	Username string `json:"username"`
}

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDetailsAll struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
	CreatedAt    string `json:"created_at"`
	RefreshToken string `json:"refresh_token"`
}

type UserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	Username string `json:"username"`
	Success  bool   `json:"success"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Username  string `json:"username"`
	AuthToken string `json:"token"`
	Success   string `json:"success"`
}
