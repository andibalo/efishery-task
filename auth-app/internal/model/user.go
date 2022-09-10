package model

type User struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	Password   string `json:"password"`
	Timestampz string `json:"timestamp"`
}
