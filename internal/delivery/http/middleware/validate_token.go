package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/usecase"
)

type ValidateTokenMiddleware struct {
	jwtKey      string
	userUsecase *usecase.UserUsecase
}

func NewValidateTokenMiddleware(
	jwtKey string,
	userUsecase *usecase.UserUsecase,
) *ValidateTokenMiddleware {
	return &ValidateTokenMiddleware{
		jwtKey:      jwtKey,
		userUsecase: userUsecase,
	}
}

func (m *ValidateTokenMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			model.ResponseError(ctx, model.ErrorUnauthorized(errors.New("authorization header is missing")))
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			model.ResponseError(ctx, model.ErrorUnauthorized(errors.New("invalid authorization format")))
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, model.ErrorUnauthorized(errors.New("unexpected signing method"))
			}
			return []byte(m.jwtKey), nil
		})
		if err != nil || !token.Valid {
			model.ResponseError(ctx, model.ErrorUnauthorized(errors.New("invalid token")))
			ctx.Abort()
			return
		}

		if jwtClaims, ok := token.Claims.(model.JWTClaims); ok {
			if _, err := m.userUsecase.GetByUsername(ctx, jwtClaims.Subject); err != nil {
				model.ResponseError(ctx, err)
				ctx.Abort()
				return
			}
			ctx.Set("jwt_claims", jwtClaims)
		}

		ctx.Next()
	}
}
