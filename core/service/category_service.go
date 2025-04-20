package service

import (
	"bwa-api/core/domain/entity"
	"bwa-api/internal/adapter/repository"
	"bwa-api/libs/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type CategoryService interface {
	GetCategory(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) (entity.CategoryEntity, error)
	EditCategory(ctx context.Context, id int64, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	slug := conv.GenerateSlug(req.Title)
	req.Slug = slug

	err := c.categoryRepository.CreateCategory(ctx, req)
	if err != nil {
		code = "[SERVICE] Create Category -1"
		log.Errorw(code, err)
		return entity.CategoryEntity{}
	}
	return nil
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(ctx context.Context, id int64) error {
	err := c.categoryRepository.DeleteCategory(ctx, id)
	if err != nil {
		code = "[SERVICE] Delete Category -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// EditCategory implements CategoryService.
func (c *categoryService) EditCategory(ctx context.Context, id int64, req entity.CategoryEntity) error {
	categoryData, err := c.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[SERVICE] Edit Category -1"
		log.Errorw(code, err)
		return err
	}
	slug := conv.GenerateSlug(req.Title)
	if categoryData.Title == req.Title {
		slug = categoryData.Slug
	}

	req.Slug = slug

	err = c.categoryRepository.EditCategory(ctx, id, req)
	if err != nil {
		code = "[SERVICE] Edit Category -2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetCategory implements CategoryService.
func (c *categoryService) GetCategory(ctx context.Context) ([]entity.CategoryEntity, error) {
	res, err := c.categoryRepository.GetCategory(ctx)
	if err != nil {
		code = "[SERVICE] Get Category -1"
		log.Errorw(code, err)
		return nil, err
	}

	return res, nil
}

// GetCategoryByID implements CategoryService.
func (c *categoryService) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	res, err := c.categoryRepository.GetCategoryByID(ctx, id)
	if err != nil {
		code = "[SERVICE] Get Category By ID -1"
		log.Errorw(code, err)
		return nil, err
	}

	return res, nil
}

// UpdateCategory implements CategoryService.
func (c *categoryService) UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) (entity.CategoryEntity, error) {
	panic("unimplemented")
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository}
}
