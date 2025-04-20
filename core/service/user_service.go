package service

import (
	"bwa-api/core/domain/entity"
	"bwa-api/internal/adapter/repository"
	"bwa-api/libs/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	UpdatePassword(ctx context.Context, newPass string, id int64) error
	GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error)
}

type userService struct {
	userRepo repository.UserRepository
}

// GetUserById implements UserService.
func (u *userService) GetUserById(ctx context.Context, id int64) (*entity.UserEntity, error) {
	res, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		code = "[USER SERVICE] Get User By ID -1"
		log.Errorw(code, err)
		return nil, err
	}

	return res, nil

}

// UpdatePassword implements UserService.
func (u *userService) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	pass, err := conv.HashPassword(newPass)
	if err != nil {
		code = "[USER SERVICE] Hash Password -1"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepo.UpdatePassword(ctx, pass, id)
	if err != nil {
		code = "[USER SERVICE] Update Password -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
