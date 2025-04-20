package service

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"bwa-api/internal/adapter/repository"
	"bwa-api/libs/auth"
	"bwa-api/libs/conv"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var code string

type AuthService interface {
	GetUserbyEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       auth.Jwt
}

// GetUserbyEmail implements AuthService.
func (a *authService) GetUserbyEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	res, err := a.authRepository.GetUserbyEmail(ctx, req)
	if err != nil {
		code = "SERVICE GetUserbyEmail -1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPassword(req.Password, res.Password); !checkPass {
		code = "[SERVICE] GetEmailByUser -2"
		err = errors.New("password not match")
		log.Errorw(code, err)
		return nil, err
	}

	JwtData := entity.JwtData{
		UserID: float64(res.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        string(res.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&JwtData)
	if err != nil {
		code = "[SERVICE] GetEmailByUser -3"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.AccessToken{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}

	return &resp, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{authRepository: authRepository, cfg: cfg, jwtToken: jwtToken}
}
