package service

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"bwa-api/internal/adapter/cloudinary"
	"bwa-api/internal/adapter/repository"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
	UploadImage(ctx context.Context, req entity.FileUploadRequest) (*entity.FileUploadResponse, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg         *config.Config
	cld         cloudinary.CloudinaryAdapter
}

// CreateContent implements ContentService.
func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err := c.contentRepo.CreateContent(ctx, req)
	if err != nil {
		code = "[CONTENT SERVICE] Create Content -1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteContent implements ContentService.
func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	contentData, err := c.contentRepo.GetContentByID(ctx, id)

	if err != nil {
		code = "[CONTENT SERVICE] Delete Content -1"
		log.Errorw(code, err)
		return err
	}

	if contentData.PublicId != "" {
		err = c.cld.DeleteImage(ctx, contentData.PublicId)
		if err != nil {
			code = "[CONTENT SERVICE] Delete Content -2"
			log.Errorw(code, err)
			return err
		}
	}

	err = c.contentRepo.DeleteContent(ctx, id)
	if err != nil {
		code = "[CONTENT SERVICE] Delete Content -3"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// EditContent implements ContentService.
func (c *contentService) EditContent(ctx context.Context, req entity.ContentEntity) error {
	contentData, err := c.contentRepo.GetContentByID(ctx, req.ID)
	if err != nil {
		code = "[CONTENT SERVICE] Edit Content -1"
		log.Errorw(code, err)
		return err
	}

	// remove old image on cloudinary only when a new image is provided
	if req.Image != "" && req.Image != contentData.Image && contentData.PublicId != "" {
		if err = c.cld.DeleteImage(ctx, contentData.PublicId); err != nil {
			code = "[CONTENT SERVICE] Edit Content -2"
			log.Errorw(code, err)
			return err
		}
	}

	// keep existing data when the request does not carry the value
	if req.Image == "" {
		req.Image = contentData.Image
	}

	if req.PublicId == "" {
		req.PublicId = contentData.PublicId
	}

	if req.CategoryId == 0 {
		req.CategoryId = contentData.CategoryId
	}

	if req.Status == "" {
		req.Status = contentData.Status
	}

	if len(req.Tags) == 0 {
		req.Tags = contentData.Tags
	}

	if req.User.ID == 0 {
		switch {
		case req.UserID != 0:
			req.User.ID = req.UserID
		case contentData.User.ID != 0:
			req.User.ID = contentData.User.ID
		}
	}

	err = c.contentRepo.UpdateContent(ctx, req)
	if err != nil {
		code = "[CONTENT SERVICE] Edit Content -3"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentService.
func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	res, err := c.contentRepo.GetContentByID(ctx, id)
	if err != nil {
		code = "[CONTENT SERVICE] Get Content By ID -1"
		log.Errorw(code, err)
		return nil, err
	}
	return res, nil
}

// GetContents implements ContentService.
func (c *contentService) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error) {
	results, totalData, totalPages, err := c.contentRepo.GetContents(ctx, query)
	if err != nil {
		code = "[CONTENT SERVICE] Get Contents -1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}
	return results, totalData, totalPages, nil
}

// UpdateContent implements ContentService.
func (c *contentService) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	// 1. Ambil data lama
	oldContent, err := c.contentRepo.GetContentByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// 2. Jika gambar berubah dan ada public ID lama, hapus gambar lama di Cloudinary
	if req.Image != "" && req.Image != oldContent.Image && oldContent.PublicId != "" {
		_ = c.cld.DeleteImage(ctx, oldContent.PublicId)
	}
	// 3. Update data content
	err = c.contentRepo.UpdateContent(ctx, req)
	if err != nil {
		code = "[CONTENT SERVICE] Update Content -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// UploadImage implements ContentService.
func (c *contentService) UploadImage(ctx context.Context, req entity.FileUploadRequest) (*entity.FileUploadResponse, error) {

	res, err := c.cld.UploadImage(ctx, &req)
	if err != nil {
		code = "[CONTENT SERVICE] Upload Image -1"
		log.Errorw(code, err)
		return nil, err
	}

	return res, nil
}

func NewContentService(contentRepo repository.ContentRepository, cfg *config.Config, cld cloudinary.CloudinaryAdapter) ContentService {
	return &contentService{contentRepo: contentRepo, cfg: cfg, cld: cld}
}
