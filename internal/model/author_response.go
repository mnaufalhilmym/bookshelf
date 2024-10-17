package model

import (
	"time"

	"github.com/mnaufalhilmym/bookshelf/internal/entity"
)

type AuthorResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
}

func ToAuthorResponse(author *entity.Author) *AuthorResponse {
	return &AuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Birthdate: author.Birthdate,
	}
}

func ToAuthorsResponse(authors []entity.Author) []AuthorResponse {
	response := make([]AuthorResponse, len(authors))
	for i, author := range authors {
		response[i] = *ToAuthorResponse(&author)
	}
	return response
}
