package model

import "github.com/mnaufalhilmym/bookshelf/internal/entity"

type BookResponse struct {
	ID         int    `json:"id"`
	Title      string `json:"title,omitempty"`
	ISBN       string `json:"isbn,omitempty"`
	AuthorID   int    `json:"author_id,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
}

func ToBookResponse(book *entity.Book) *BookResponse {
	return &BookResponse{
		ID:         book.ID,
		Title:      book.Title,
		ISBN:       book.ISBN,
		AuthorID:   book.AuthorID,
		AuthorName: book.Author.Name,
	}
}

func ToBooksResponse(books []entity.Book) []BookResponse {
	response := make([]BookResponse, len(books))
	for i, book := range books {
		response[i] = *ToBookResponse(&book)
	}
	return response
}