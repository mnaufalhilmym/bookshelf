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

	v1.POST("/auth/register", r.userHandler.Register)
	v1.POST("/auth/login", r.userHandler.Login)

	v1.GET("/authors", r.authorHandler.GetMany)
	v1.GET("/authors/:id", r.authorHandler.Get)
	v1.POST("/authors", r.authorHandler.Create)
	v1.PUT("/authors/:id", r.authorHandler.Update)
	v1.DELETE("/authors/:id", r.authorHandler.Delete)

	v1.GET("/books", r.authorHandler.GetMany)
	v1.GET("/books/:id", r.authorHandler.Get)
	v1.POST("/books", r.authorHandler.Create)
	v1.PUT("/books/:id", r.authorHandler.Update)
	v1.DELETE("/books/:id", r.authorHandler.Delete)
}
