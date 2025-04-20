package handler

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/service"
	"bwa-api/internal/adapter/handler/request"
	"bwa-api/internal/adapter/handler/response"
	"bwa-api/libs/conv"
	"bwa-api/libs/validator"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImage(c *fiber.Ctx) error

	GetContentWithQuery(c *fiber.Ctx) error
	GetContentDetail(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// GetContentDetail implements ContentHandler.
func (ch *contentHandler) GetContentDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	contenId, err := conv.StringToInt64(id)
	if err != nil {
		code = "[CONTENT HANDLER] GetContentDetail -1"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	res, err := ch.contentService.GetContentByID(c.Context(), contenId)
	if err != nil {
		code = "[CONTENT HANDLER] GetContentDetail -2"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp := response.ContentResponse{
		ID:          res.ID,
		Title:       res.Title,
		Excerpt:     res.Excerpt,
		Description: res.Description,
		Image:       res.Image,
		Tags:        res.Tags,
		Status:      res.Status,
		CategoryId:  res.CategoryId,
		UserID:      res.User.ID,
		CreatedAt:   res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Format("2006-01-02 15:04:05"),
		User:        res.User.Name,
		Category:    res.Category.Title,
	}

	defaultResponse.Data = resp
	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetContentWithQuery implements ContentHandler.
func (ch *contentHandler) GetContentWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[CONTENT HANDLER] GetContentWithQuery -1"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[CONTENT HANDLER] GetContentWithQuery -2"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	orderBy := "created_at"
	if c.Query("order_by") != "" {
		orderBy = c.Query("order_by")
	}

	orderType := "desc"
	if c.Query("order_type") != "" {
		orderType = c.Query("order_type")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	CategoryId := 0
	if c.Query("category_id") != "" {
		CategoryId, err = conv.StringToInt(c.Query("category_id"))
		if err != nil {
			code = "[CONTENT HANDLER] GetContentWithQuery -3"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}
	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		Search:     search,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Status:     "Published",
		CategoryId: int64(CategoryId),
	}
	// UserId := claims.UserID
	result, totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[CONTENT HANDLER] GetContentWithQuery -4"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Meta.Status = true

	respContent := []response.ContentResponse{}
	for _, val := range result {
		resp := response.ContentResponse{
			ID:          val.ID,
			Title:       val.Title,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			Tags:        val.Tags,
			Status:      val.Status,
			CategoryId:  val.CategoryId,
			UserID:      val.User.ID,
			CreatedAt:   val.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   val.UpdatedAt.Format("2006-01-02 15:04:05"),
			User:        val.User.Name,
			Category:    val.Category.Title,
		}

		respContent = append(respContent, resp)

	}

	defaultResponse.Data = respContent
	defaultResponse.Pagina = &response.PaginationResponse{
		TotalRecords: totalData,
		TotalPages:   totalPages,
		Page:         int64(page),
		PerPage:      int64(limit),
	}
	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// CreateContent implements ContentHandler.
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Create  Content -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}
	// userId := claims.UserID
	var req request.CreateContent
	if err := c.BodyParser(&req); err != nil {
		code = "[CONTENT HANDLER] Create Content - 2"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}
	userId := claims.UserID

	err = validator.ValidateStruct(req)
	if err != nil {
		code = "[CONTENT HANDLER] Create Content - 3"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryId:  int64(req.CategoryID),
		UserID:      int64(userId),
	}

	err := ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[CONTENT HANDLER] Create Content - 4"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Message = "Success"
	defaultResponse.Meta.Status = true
	defaultResponse.Data = nil
	defaultResponse.Pagina = nil

	return c.Status(fiber.StatusCreated).JSON(defaultResponse)
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Delete Content -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	id := c.Params("id")
	contendId, err := conv.StringToInt64(id)
	if err != nil {
		code = "[CONTENT HANDLER] Delete Content -2"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), contendId)
	if err != nil {
		code = "[CONTENT HANDLER] Delete Content -3"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Meta.Status = true
	defaultResponse.Data = nil
	defaultResponse.Pagina = nil

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetContentByID implements ContentHandler.
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Get Contents by ID -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	id := c.Params("id")
	contenId, err := conv.StringToInt64(id)
	if err != nil {
		code = "[CONTENT HANDLER] Get Contents by ID -2"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	res, err := ch.contentService.GetContentByID(c.Context(), contenId)
	if err != nil {
		code = "[CONTENT HANDLER] Get Contents by ID -3"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp := response.ContentResponse{
		ID:          res.ID,
		Title:       res.Title,
		Excerpt:     res.Excerpt,
		Description: res.Description,
		Image:       res.Image,
		Tags:        res.Tags,
		Status:      res.Status,
		CategoryId:  res.CategoryId,
		UserID:      res.User.ID,
		CreatedAt:   res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   res.UpdatedAt.Format("2006-01-02 15:04:05"),
		User:        res.User.Name,
		Category:    res.Category.Title,
	}

	defaultResponse.Data = resp
	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetContents implements ContentHandler.
func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Get Contents -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[CONTENT HANDLER] Get Contents -2"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[CONTENT HANDLER] Get Contents -3"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	orderBy := "created_at"
	if c.Query("order_by") != "" {
		orderBy = c.Query("order_by")
	}

	orderType := "desc"
	if c.Query("order_type") != "" {
		orderType = c.Query("order_type")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	CategoryId := 0
	if c.Query("category_id") != "" {
		CategoryId, err = conv.StringToInt(c.Query("category_id"))
		if err != nil {
			code = "[CONTENT HANDLER] Get Contents -4"
			log.Errorw(code, err)
			errResp.Message = err.Error()
			errResp.Status = false
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		Search:     search,
		OrderBy:    orderBy,
		OrderType:  orderType,
		CategoryId: int64(CategoryId),
		Status:     "Published",
	}
	// UserId := claims.UserID
	result, totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[CONTENT HANDLER] Get Contents -5"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Meta.Status = true

	respContent := []response.ContentResponse{}
	for _, val := range result {
		resp := response.ContentResponse{
			ID:          val.ID,
			Title:       val.Title,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			Tags:        val.Tags,
			Status:      val.Status,
			CategoryId:  val.CategoryId,
			UserID:      val.User.ID,
			CreatedAt:   val.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   val.UpdatedAt.Format("2006-01-02 15:04:05"),
			User:        val.User.Name,
			Category:    val.Category.Title,
		}

		respContent = append(respContent, resp)

	}

	defaultResponse.Data = respContent
	defaultResponse.Pagina = &response.PaginationResponse{
		TotalRecords: totalData,
		TotalPages:   totalPages,
		Page:         int64(page),
		PerPage:      int64(limit),
	}
	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// UpdateContent implements ContentHandler.
func (ch *contentHandler) UpdateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Update  Content -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}
	var req request.CreateContent
	if err := c.BodyParser(&req); err != nil {
		code = "[CONTENT HANDLER] Update Content - 2"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	userId := claims.UserID

	err = validator.ValidateStruct(req)
	if err != nil {
		code = "[CONTENT HANDLER] Update Content - 3"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	id := c.Params("id")
	contendId, err := conv.StringToInt64(id)
	if err != nil {
		code = "[CONTENT HANDLER] Update Content - 4"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}
	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          contendId,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryId:  int64(req.CategoryID),
		UserID:      int64(userId),
	}

	err = ch.contentService.UpdateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[CONTENT HANDLER] Update Content - 5"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	defaultResponse.Meta.Message = "Success"
	defaultResponse.Meta.Status = true
	defaultResponse.Data = nil

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// UploadImage implements ContentHandler.
func (ch *contentHandler) UploadImage(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[CONTENT HANDLER] Upload Image -1"
		log.Errorw(code, "user not found", claims)
		errResp.Message = "Unauthorized"
		errResp.Status = false
		return c.Status(fiber.StatusUnauthorized).JSON(errResp)
	}

	var req request.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code = "[CONTENT HANDLER] Upload Image -2"
		log.Errorw(code, err)
		errResp.Message = err.Error()
		errResp.Status = false
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	if err := c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code = "[CONTENT HANDLER] Upload Image -3"
		log.Errorw(code, err)
		errResp.Meta.Message = "Failed to save file"
		errResp.Meta.Status = false

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadRequest{
		Name: fmt.Sprintf("%d-%d", int64(claims.UserID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imgUrl, err := ch.contentService.UploadImage(c.Context(), reqEntity)
	if err != nil {
		code = "[CONTENT HANDLER] Upload Image -4"
		log.Errorw(code, err)
		errResp.Meta.Message = err.Error()
		errResp.Meta.Status = false

		return c.Status(fiber.StatusInternalServerError).JSON(errResp)

	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code = "[CONTENT HANDLER] Upload Image -5"
			log.Errorw(code, err)
			errResp.Meta.Message = err.Error()
			errResp.Meta.Status = false

			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
	}

	imageUrlResp := map[string]interface{}{
		"urlImage": imgUrl,
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "Success"
	defaultResponse.Data = imageUrlResp
	defaultResponse.Pagina = nil

	return c.Status(fiber.StatusCreated).JSON(defaultResponse)

}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
