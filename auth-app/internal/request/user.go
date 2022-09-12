package request

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

type LoginUserRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetUserTokenDetailsRequest struct {
	Token string `json:"token" validate:"required"`
}
