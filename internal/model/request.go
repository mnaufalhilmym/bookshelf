package model

type paginationRequest struct {
	Page int `form:"page" binding:"omitempty,gt=0"`
	Size int `form:"size" binding:"omitempty,gt=0"`
}
