package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mnaufalhilmym/bookshelf/internal/entity"
	"github.com/mnaufalhilmym/bookshelf/internal/model"
	"github.com/mnaufalhilmym/bookshelf/internal/repository"
	"github.com/mnaufalhilmym/gotracing"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase struct {
	db            *gorm.DB
	repository    *repository.UserRepository
	jwtKey        string
	jwtExpiration time.Duration
}

func NewUserUsecase(
	db *gorm.DB,
	repository *repository.UserRepository,
	jwtKey string,
	jwtExpiration time.Duration,
) *UserUsecase {
	return &UserUsecase{
		db,
		repository,
		jwtKey,
		jwtExpiration,
	}
}

func (uc *UserUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := uc.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		gotracing.Error("Failed to generate hashed password", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to generate hashed password"))
	}

	user := &entity.User{
		Username: request.Username,
		Password: string(hashedPassword),
	}

	if err := uc.repository.Create(tx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, model.ErrorBadRequest(errors.New("duplicate username"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to create new user"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToUserResponse(user), nil
}

func (uc *UserUsecase) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	user, err := uc.repository.FindByUsername(tx, request.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("username not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find user data by username"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, model.ErrorBadRequest(errors.New("wrong password"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to compare hash and password"))
	}

	jwtClaims := model.JWTClaims{ID: user.ID}
	jwtClaims.Issuer = "bookshelf-server"
	jwtClaims.Subject = user.Username
	jwtClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(uc.jwtExpiration))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(uc.jwtKey))
	if err != nil {
		gotracing.Error("Failed to sign JWT token", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to sign JWT token"))
	}

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	return model.ToLoginResponse(user, tokenString), nil
}

func (uc *UserUsecase) GetByUsername(ctx context.Context, username string) (*model.UserResponse, error) {
	tx := uc.db.WithContext(ctx).Begin(&sql.TxOptions{ReadOnly: true})
	defer tx.Rollback()

	if err := tx.Commit().Error; err != nil {
		gotracing.Error("Failed to commit transaction", err)
		return nil, model.ErrorInternalServerError(errors.New("failed to commit transaction"))
	}

	user, err := uc.repository.FindByUsername(tx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrorNotFound(errors.New("username not found"))
		}
		return nil, model.ErrorInternalServerError(errors.New("failed to find user data by username"))
	}

	return model.ToUserResponse(user), nil
}
