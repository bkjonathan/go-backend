package auth

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type TokenResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}
