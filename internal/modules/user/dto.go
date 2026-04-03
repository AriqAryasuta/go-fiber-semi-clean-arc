package user

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=120"`
	Email string `json:"email" validate:"required,email,max=160"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
