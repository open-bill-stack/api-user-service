package structure

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RefreshTokenRequest struct {
	Token string `json:"token"`
}
