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
	r.router.POST("/auth/register", r.userHandler.Register)
	r.router.POST("/auth/login", r.userHandler.Login)

	r.router.Use(r.validateTokenMiddleware.ValidateToken())

	r.router.GET("/authors", r.authorHandler.GetMany)
	r.router.GET("/authors/:id", r.authorHandler.Get)
	r.router.POST("/authors", r.authorHandler.Create)
	r.router.PUT("/authors/:id", r.authorHandler.Update)
	r.router.DELETE("/authors/:id", r.authorHandler.Delete)

	r.router.GET("/books", r.bookHandler.GetMany)
	r.router.GET("/books/:id", r.bookHandler.Get)
	r.router.POST("/books", r.bookHandler.Create)
	r.router.PUT("/books/:id", r.bookHandler.Update)
	r.router.DELETE("/books/:id", r.bookHandler.Delete)
}
