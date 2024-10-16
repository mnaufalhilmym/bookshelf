package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/handler"
)

type RouteConfig struct {
	router        *gin.Engine
	userHandler   *handler.UserHandler
	authorHandler *handler.AuthorHandler
	bookHandler   *handler.BookHandler
}

func New(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	authorHandler *handler.AuthorHandler,
	bookHandler *handler.BookHandler,
) *RouteConfig {
	return &RouteConfig{
		router,
		userHandler,
		authorHandler,
		bookHandler,
	}
}

func (r *RouteConfig) ConfigureRoutes() {
	v1 := r.router.Group("/v1")

	v1.POST("/auth/register")
	v1.POST("/auth/login")

	v1.GET("/authors")
	v1.GET("/authors/:id")
	v1.POST("/authors")
	v1.PUT("/authors/:id")
	v1.DELETE("/authors/:id")

	v1.GET("/books")
	v1.GET("/books/:id")
	v1.POST("/books")
	v1.PUT("/books/:id")
	v1.DELETE("/books/:id")
}
