package repository

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/domain/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepository interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userRepository struct {
	db *gorm.DB
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User
	err = u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code = "[USER REPOSITORY] Get User By ID -1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:    modelUser.ID,
		Name:  modelUser.Name,
		Email: modelUser.Email,
	}, nil
}

// UpdatePassword implements UserRepository.
func (u *userRepository) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code = "[USER REPOSITORY] Update Password -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
