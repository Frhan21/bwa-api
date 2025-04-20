package service

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"bwa-api/internal/adapter/cloudflare"
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
	UploadImage(ctx context.Context, req entity.FileUploadRequest) (string, error)
}

type contentService struct {
	contentRepo repository.ContentRepository
	cfg         *config.Config
	r2          cloudflare.CloudFlareR2Adapter
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
	err := c.contentRepo.DeleteContent(ctx, id)
	if err != nil {
		code = "[CONTENT SERVICE] Delete Content -1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// EditContent implements ContentService.
func (c *contentService) EditContent(ctx context.Context, req entity.ContentEntity) error {

	panic("unimplemented")
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
	err := c.contentRepo.UpdateContent(ctx, req)
	if err != nil {
		code = "[CONTENT SERVICE] Update Content -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// UploadImage implements ContentService.
func (c *contentService) UploadImage(ctx context.Context, req entity.FileUploadRequest) (string, error) {

	// Upload to R2
	url, err := c.r2.UploadImage(&req)
	if err != nil {
		code = "[CONTENT SERVICE] Upload Image -1"
		log.Errorw(code, err)
		return "", err
	}

	return url, nil
}

func NewContentService(contentRepo repository.ContentRepository, cfg *config.Config, r2 cloudflare.CloudFlareR2Adapter) ContentService {
	return &contentService{contentRepo: contentRepo, cfg: cfg, r2: r2}
}
