package handler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
	"github.com/mnaufalhilmym/gotracing"
)

type BookHandler struct {
	usecase *usecase.BookUsecase
}

func NewBookHandler(uc *usecase.BookUsecase) *BookHandler {
	return &BookHandler{uc}
}

func (h *BookHandler) GetMany(ctx *gin.Context) {
	request := new(model.GetManyBooksRequest)
	if err := ctx.ShouldBindQuery(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}

	response, total, err := h.usecase.GetMany(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOKPaginated(ctx, response, total, request.Page, request.Size)
}

func (h *BookHandler) Get(ctx *gin.Context) {
	request := new(model.GetBookRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Get(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, response)
}

func (h *BookHandler) Create(ctx *gin.Context) {
	request := new(model.CreateBookRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Create(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseCreated(ctx, response)
}

func (h *BookHandler) Update(ctx *gin.Context) {
	request := new(model.UpdateBookRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Update(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, response)
}

func (h *BookHandler) Delete(ctx *gin.Context) {
	request := new(model.DeleteBookRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.ErrorBadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.ErrorBadRequest(errors.New("failed to parse request")))
		return
	}

	bookID, err := h.usecase.Delete(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, *bookID)
}
