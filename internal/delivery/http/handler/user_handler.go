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

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	request := new(model.RegisterUserRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.BadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Register(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseCreated(ctx, response)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	request := new(model.LoginRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		gotracing.Error("Failed to parse request", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			model.ResponseError(ctx, model.BadRequest(fmt.Errorf("validation error in field %s", errs[0].Field())))
			return
		}
		model.ResponseError(ctx, model.BadRequest(errors.New("failed to parse request")))
		return
	}

	response, err := h.usecase.Login(ctx, request)
	if err != nil {
		model.ResponseError(ctx, err)
		return
	}

	model.ResponseOK(ctx, response)
}
