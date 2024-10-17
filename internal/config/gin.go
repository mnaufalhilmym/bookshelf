package config

import (
	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
)

func NewGin(appMode string) *gin.Engine {
	gin.SetMode(appMode)

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(customErrorHandler())

	return router
}

func customErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			model.ResponseError(ctx, model.ErrorInternalServerError(ctx.Errors[0]))
		}
	}
}
