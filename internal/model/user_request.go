package model

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,gt=0"`
	Password string `json:"password" binding:"required,gt=0"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,gt=0"`
	Password string `json:"password" binding:"required,gt=0"`
}
