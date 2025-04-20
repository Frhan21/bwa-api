package handler

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/service"
	"bwa-api/internal/adapter/handler/request"
	"bwa-api/internal/adapter/handler/response"
	"bwa-api/libs/conv"
	"bwa-api/libs/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultResponse response.DefaultSuccessResponse

type CategoryHandler interface {
	GetCategory(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	CreateCategory(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
	EditCategory(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error

	GetCategoryFE(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// GetCategoryFE implements CategoryHandler.
func (ch *categoryHandler) GetCategoryFE(c *fiber.Ctx) error {
	results, err := ch.categoryService.GetCategory(c.Context())
	if err != nil {
		code := "[HANDLER] Get Category FE -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to get category"
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, val := range results {
		category := response.SuccessCategoryResponse{
			ID:    int64(val.ID),
			Title: val.Title,
			Slug:  val.Slug,
		}

		if val.User.ID != 0 { // Check for zero-value struct
			category.User = response.UserResponse{
				ID:    val.User.ID,
				Name:  val.User.Name,
				Email: val.User.Email,
			}
		}
		categoryResponses = append(categoryResponses, category)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = categoryResponses

	return c.Status(fiber.StatusOK).JSON(defaultResponse)

}

// CreateCategory implements CategoryHandler.
func (c *categoryHandler) CreateCategory(ctx *fiber.Ctx) error {
	var req request.CreateCategoryRequest
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] Create Category -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"
		return ctx.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	if err = ctx.BodyParser(&req); err != nil {
		code = "[HANDLER] Create Category -2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to parse request"
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validator.ValidateStruct(req); err != nil {
		code = "[HANDLER] Create Category -3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User:  entity.UserEntity{ID: int64(userID)},
	}

	err = c.categoryService.CreateCategory(ctx.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] Create Category -4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to create category"
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = nil

	return ctx.Status(fiber.StatusCreated).JSON(defaultResponse)
}

// DeleteCategory implements CategoryHandler.
func (c *categoryHandler) DeleteCategory(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] Delete Category -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"
		return ctx.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParams := ctx.Params("id")
	id, err := conv.StringToInt64(idParams)
	if err != nil {
		code = "[HANDLER] Delete Category -2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	err = c.categoryService.DeleteCategory(ctx.Context(), id)
	if err != nil {
		code = "[HANDLER] Delete Category -3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to delete category"
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = nil

	return ctx.Status(fiber.StatusOK).JSON(defaultResponse)
}

// EditCategory implements CategoryHandler.
func (c *categoryHandler) EditCategory(ctx *fiber.Ctx) error {
	var req request.CreateCategoryRequest
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] EditCategory -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"
		return ctx.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	if err = ctx.BodyParser(&req); err != nil {
		code = "[HANDLER] EditCategory -2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to parse request"
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err = validator.ValidateStruct(req); err != nil {
		code = "[HANDLER] EditCategory -3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	idParams := ctx.Params("id")
	id, err := conv.StringToInt64(idParams)
	if err != nil {
		code = "[HANDLER] EditCategory By ID -4"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	reqEntity := entity.CategoryEntity{
		ID:    int(id),
		Title: req.Title,
		User:  entity.UserEntity{ID: int64(userID)},
	}

	err = c.categoryService.EditCategory(ctx.Context(), id, reqEntity)
	if err != nil {
		code = "[HANDLER] EditCategory -5"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to edit category"
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = nil

	return ctx.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetCategory implements CategoryHandler.
func (c *categoryHandler) GetCategory(ctx *fiber.Ctx) error {
	var errResp response.ErrorResponseDefault
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] Get Category -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"
		return ctx.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	res, err := c.categoryService.GetCategory(ctx.Context())
	if err != nil {
		code = "[HANDLER] Get Category -2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to get category"
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResp := []response.SuccessCategoryResponse{}
	if len(res) == 0 {
		log.Warn("[HANDLER] Get Category - No categories found")
		return ctx.JSON(response.DefaultSuccessResponse{
			Meta: response.Meta{
				Status:  true,
				Message: "Success",
			},
			Data: []response.SuccessCategoryResponse{}, // Return an empty array
		})
	}

	for _, val := range res {
		category := response.SuccessCategoryResponse{
			ID:    int64(val.ID),
			Title: val.Title,
			Slug:  val.Slug,
		}

		if val.User.ID != 0 { // Check for zero-value struct
			category.User = response.UserResponse{
				ID:    val.User.ID,
				Name:  val.User.Name,
				Email: val.User.Email,
			}
		}
		categoryResp = append(categoryResp, category)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = categoryResp

	return ctx.JSON(defaultResponse)
}

// GetCategoryByID implements CategoryHandler.
func (c *categoryHandler) GetCategoryByID(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] Get Category By ID -1"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Unauthorized"
		return ctx.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	idParams := ctx.Params("id")
	id, err := conv.StringToInt64(idParams)
	if err != nil {
		code = "[HANDLER] Get Category By ID -2"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = err.Error()
		return ctx.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	res, err := c.categoryService.GetCategoryByID(ctx.Context(), id)
	if err != nil {
		code = "[HANDLER] Get Category By ID -3"
		log.Errorw(code, err)
		errResp.Meta.Status = false
		errResp.Meta.Message = "Failed to get category"
		return ctx.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	categoryResp := response.SuccessCategoryResponse{
		ID:    int64(res.ID),
		Title: res.Title,
		Slug:  res.Slug,
		User: response.UserResponse{
			ID:    res.User.ID,
			Name:  res.User.Name,
			Email: res.User.Email,
		},
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Pagina = nil
	defaultResponse.Data = categoryResp

	return ctx.JSON(defaultResponse)
}

// UpdateCategory implements CategoryHandler.
func (c *categoryHandler) UpdateCategory(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
