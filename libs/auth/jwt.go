package auth

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(data *entity.JwtData) (string, int64, error)
	VerifyAccessToken(token string) (*entity.JwtData, error)
}

type Options struct {
	signingKey string
	issuer     string
}

func (o *Options) GenerateToken(data *entity.JwtData) (string, int64, error) {
	now := time.Now().Local()
	expiredAt := now.Add(time.Hour + 24)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiredAt)
	data.RegisteredClaims.Issuer = o.issuer
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	accessTokenString, err := accessToken.SignedString([]byte(o.signingKey))
	if err != nil {
		return "", 0, err
	}
	return accessTokenString, expiredAt.Unix(), nil
}

func (o *Options) VerifyAccessToken(tokenString string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(o.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	jwtData := entity.JwtData{
		UserID: claims["user_id"].(float64),
	}
	return &jwtData, nil
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.signingKey = cfg.App.JwtSecretKey
	opt.issuer = cfg.App.JwtIssuer
	return opt
}
