package model

import "github.com/mnaufalhilmym/bookshelf/internal/entity"

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}

type LoginResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func ToLoginResponse(user *entity.User, token string) *LoginResponse {
	return &LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    token,
	}
}
