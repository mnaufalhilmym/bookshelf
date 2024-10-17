package model

import "time"

type GetManyAuthorsRequest struct {
	paginationRequest
	Name           *string    `form:"name"` // case insensitive | contains
	BirthdateStart *time.Time `form:"birthdate_start"`
	BirthdateEnd   *time.Time `form:"birthdate_end"`
}

type GetAuthorRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}

type CreateAuthorRequest struct {
	Name      string    `json:"name" binding:"required,gt=0"`
	Birthdate time.Time `json:"birthdate" binding:"required"`
}

type UpdateAuthorRequest struct {
	ID        int        `json:"-" uri:"id" binding:"required,gt=0"`
	Name      *string    `json:"name" uri:"-"`
	Birthdate *time.Time `json:"birthdate" uri:"-"`
}

type DeleteAuthorRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}
