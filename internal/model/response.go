package model

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/util"
)

type Response[T any] struct {
	Error      string      `json:"error,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
	Data       T           `json:"data"`
	DataHash   string      `json:"data_hash"`
}

type pagination struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func ResponseCreated[T any](ctx *gin.Context, data T) {
	dataJSON, _ := json.Marshal(data)
	ctx.JSON(http.StatusCreated, Response[T]{
		Data:     data,
		DataHash: util.CalculateHash(string(dataJSON)),
	})
}

func ResponseOK[T any](ctx *gin.Context, data T) {
	dataJSON, _ := json.Marshal(data)
	ctx.JSON(http.StatusOK, Response[T]{
		Data:     data,
		DataHash: util.CalculateHash(string(dataJSON)),
	})
}

func ResponseOKPaginated[T any](ctx *gin.Context, data []T, totalItem int64, page int, size int) {
	dataJSON, _ := json.Marshal(data)
	ctx.JSON(http.StatusOK, Response[[]T]{
		Data:     data,
		DataHash: util.CalculateHash(string(dataJSON)),
		Pagination: &pagination{
			Page:      page,
			Size:      size,
			TotalItem: totalItem,
			TotalPage: int64(math.Ceil(float64(totalItem) / float64(size))),
		},
	})
}

func ResponseError(ctx *gin.Context, err error) {
	appError, ok := err.(*Error)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Response[any]{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(appError.Code, Response[any]{
		Error: appError.Err.Error(),
	})
}
