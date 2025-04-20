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

type UserHandler interface {
	GetUserById(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserById implements UserHandler.
func (u *userHandler) GetUserById(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[USER HANDLER] Get User By ID -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return ctx.Status(fiber.StatusNotFound).JSON(errResp)
	}

	user, err := u.userService.GetUserById(ctx.Context(), int64(claims.UserID))
	if err != nil {
		code = "[USER HANDLER] Get User By ID -2"
		log.Errorw(code, err)
		errResp.Message = "User not found"
		errResp.Status = false
		return ctx.Status(fiber.StatusNotFound).JSON(errResp)
	}
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	defaultResponse.Data = resp
	return ctx.Status(fiber.StatusOK).JSON(defaultResponse)
}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[USER HANDLER] Update Password -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return ctx.Status(fiber.StatusNotFound).JSON(errResp)
	}

	var req request.UpdatePassword
	if err = ctx.BodyParser(&req); err != nil {
		code = "[USER HANDLER] Update Password -2"
		log.Errorw(code, err)
		errResp.Message = "Invalid request"
		errResp.Status = false
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validator.ValidateStruct(req); err != nil {
		code = "[USER HANDLER] Update Password -3"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = u.userService.UpdatePassword(ctx.Context(), req.NewPassword, int64(claims.UserID))
	if err != nil {
		code = "[USER HANDLER] Update Password -4"
		log.Errorw(code, err)
		errResp.Message = "Failed to update password"
		errResp.Status = false
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Data = nil
	return ctx.Status(fiber.StatusOK).JSON(defaultResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}
