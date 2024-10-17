package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/handler"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/middleware"
)

type RouteConfig struct {
	router *gin.Engine

	userHandler   *handler.UserHandler
	authorHandler *handler.AuthorHandler
	bookHandler   *handler.BookHandler

	validateTokenMiddleware *middleware.ValidateTokenMiddleware
}

func New(
	router *gin.Engine,

	userHandler *handler.UserHandler,
	authorHandler *handler.AuthorHandler,
	bookHandler *handler.BookHandler,

	validateTokenMiddleware *middleware.ValidateTokenMiddleware,
) *RouteConfig {
	return &RouteConfig{
		router,
		userHandler,
		authorHandler,
		bookHandler,
		validateTokenMiddleware,
	}
}

func (r *RouteConfig) ConfigureRoutes() {
	v1 := r.router.Group("/v1")

	v1.POST("/auth/register", r.userHandler.Register)
	v1.POST("/auth/login", r.userHandler.Login)

	v1.Use(r.validateTokenMiddleware.ValidateToken())

	v1.GET("/authors", r.authorHandler.GetMany)
	v1.GET("/authors/:id", r.authorHandler.Get)
	v1.POST("/authors", r.authorHandler.Create)
	v1.PUT("/authors/:id", r.authorHandler.Update)
	v1.DELETE("/authors/:id", r.authorHandler.Delete)

	v1.GET("/books", r.bookHandler.GetMany)
	v1.GET("/books/:id", r.bookHandler.Get)
	v1.POST("/books", r.bookHandler.Create)
	v1.PUT("/books/:id", r.bookHandler.Update)
	v1.DELETE("/books/:id", r.bookHandler.Delete)
}
