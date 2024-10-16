package config

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/handler"
	"github.com/mnaufalhilmym/bookshelf/internal/delivery/http/route"
	"github.com/mnaufalhilmym/bookshelf/internal/repository"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
	"gorm.io/gorm"
)

func Bootstrap(
	router *gin.Engine,
	db *gorm.DB,
	jwtKey string,
	jwtExpiration time.Duration,
) {
	// Repository
	userRepository := repository.NewUserRepository(db)
	authorRepository := repository.NewAuthorRepository(db)
	bookRepository := repository.NewBookRepository(db)

	// Usecase
	userUsecase := usecase.NewUserUsecase(db, userRepository, jwtKey, jwtExpiration)
	authorUsecase := usecase.NewAuthorUsecase(db, authorRepository)
	bookUsecase := usecase.NewBookUsecase(db, bookRepository, authorRepository)

	// Handler
	userHandler := handler.NewUserHandler(userUsecase)
	authorHandler := handler.NewAuthorHandler(authorUsecase)
	bookHandler := handler.NewBookHandler(bookUsecase)

	routeConfig := route.New(
		router,
		userHandler,
		authorHandler,
		bookHandler,
	)

	routeConfig.ConfigureRoutes()
}
