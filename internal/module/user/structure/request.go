package structure

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserRequest struct {
	UUID string `json:"uuid" validate:"required"`
}
