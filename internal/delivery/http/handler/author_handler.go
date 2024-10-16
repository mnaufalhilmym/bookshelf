package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
	"github.com/mnaufalhilmym/gotracing"
)

type AuthorHandler struct {
	usecase *usecase.AuthorUsecase
}

func NewAuthorHandler(uc *usecase.AuthorUsecase) *AuthorHandler {
	return &AuthorHandler{uc}
}

func (h *AuthorHandler) GetMany(ctx *gin.Context) {
	request := new(model.GetManyAuthorsRequest)
	if err := ctx.ShouldBindQuery(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
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

func (h *AuthorHandler) Get(ctx *gin.Context) {
	request := new(model.GetAuthorRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Get(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, response)
}

func (h *AuthorHandler) Create(ctx *gin.Context) {
	request := new(model.CreateAuthorRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Create(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseCreated(ctx, response)
}

func (h *AuthorHandler) Update(ctx *gin.Context) {
	request := new(model.UpdateAuthorRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Update(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, response)
}

func (h *AuthorHandler) Delete(ctx *gin.Context) {
	request := new(model.DeleteAuthorRequest)
	if err := ctx.ShouldBindUri(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	authorID, err := h.usecase.Delete(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, model.AuthorResponse{ID: *authorID})
}
