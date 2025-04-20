package handler

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/service"
	"bwa-api/internal/adapter/handler/request"
	"bwa-api/internal/adapter/handler/response"
	"bwa-api/libs/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errResp response.ErrorResponseDefault

// var val = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHanlder.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := request.LoginRequest{}
	res := response.SucccessLoginResponse{}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] Login 1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validator.ValidateStruct(req); err != nil {
		code = "[HANDL	ER] Login 2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqLogin := entity.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserbyEmail(c.Context(), reqLogin)
	if err != nil {
		code = "[HANDLER] Login 3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()

		if err.Error() == "record not found" {
			errResp.Meta.Status = false
			errResp.Meta.Message = "password not match"
			return c.Status(fiber.StatusUnauthorized).JSON(errResp)
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	res.Meta.Status = true
	res.Meta.Message = "success"
	res.ExpiresAt = result.ExpiresAt
	res.AccessToken = result.AccessToken

	return c.JSON(res)
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{authService: authService}
}
