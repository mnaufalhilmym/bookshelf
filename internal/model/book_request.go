package model

type GetManyBooksRequest struct {
	pagination
	Title      *string `form:"title"` // case insensitive | contains
	ISBN       *string `form:"isbn"`  // case insensitive | contains
	AuthorID   *int    `form:"author_id" binding:"omitempty,gt=0"`
	AuthorName *string `form:"author_name"` // case insensitive | contains
}

type GetBookRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}

type CreateBookRequest struct {
	Title    string `json:"title" binding:"required"`
	ISBN     string `json:"isbn" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required,gt=0"`
}

type UpdateBookRequest struct {
	ID       int     `json:"-" uri:"id" binding:"required,gt=0"`
	Title    *string `json:"title" uri:"-"`
	ISBN     *string `json:"isbn" uri:"-"`
	AuthorID *int    `json:"author_id" binding:"omitempty,gt=0"`
}

type DeleteBookRequest struct {
	ID int `uri:"id" binding:"required,gt=0"`
}
