package request

type CreateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type LoginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
