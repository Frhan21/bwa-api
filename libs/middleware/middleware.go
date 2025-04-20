package middleware

import (
	"bwa-api/config"
	"bwa-api/internal/adapter/handler/response"
	"bwa-api/libs/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	CheckToken() fiber.Handler
}

type Options struct {
	authJwt auth.Jwt
}

func (o *Options) CheckToken() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var errorResponse response.ErrorResponseDefault
		authHandler := c.Get("Authorization")
		if authHandler == " " {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Unauthorized"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		tokenString := strings.Split(authHandler, "Bearer ")[1]
		claims, err := o.authJwt.VerifyAccessToken(tokenString)
		if err != nil {
			errorResponse.Meta.Status = false
			errorResponse.Meta.Message = "Unauthorized"
			return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
		}

		c.Locals("user", claims)
		return c.Next()
	}

}

func NewMiddleware(cfg *config.Config) *Options {
	opt := new(Options)
	opt.authJwt = auth.NewJwt(cfg)
	return opt
}
