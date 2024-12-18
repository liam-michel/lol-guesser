package common

type Username struct {
	Username string `json:"username"`
}

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
