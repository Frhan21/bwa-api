package repository

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/domain/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var err error
var code string

type AuthRepository interface {
	GetUserbyEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

func (a *authRepository) GetUserbyEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error) {
	var modelUser model.User
	err = a.db.Where("email = ?", req.Email).First(&modelUser).Error
	if err != nil {
		code = "REPOSITORY GetUserbyEmail -1"
		log.Errorw(code, err)
		return nil, err
	}

	res := entity.UserEntity{
		ID:       modelUser.ID,
		Email:    modelUser.Email,
		Name:     modelUser.Name,
		Password: modelUser.Password,
	}

	return &res, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
