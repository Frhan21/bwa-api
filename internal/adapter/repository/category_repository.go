package repository

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/domain/model"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategory(ctx context.Context) ([]entity.CategoryEntity, error)
	GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error)
	CreateCategory(ctx context.Context, req entity.CategoryEntity) error
	UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) (entity.CategoryEntity, error)
	EditCategory(ctx context.Context, id int64, req entity.CategoryEntity) error
	DeleteCategory(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *gorm.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, req entity.CategoryEntity) error {
	var countSlug int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&countSlug).Error
	if err != nil {
		code = "[REPOSITORY] Create Category -2"
		log.Errorw(code, err)
		return err
	}

	countSlug = countSlug + 1
	slug := fmt.Sprintf("%s-%d", req.Slug, countSlug)
	modelCategory := model.Category{
		Title: req.Title,
		Slug:  slug,
		User:  model.User{ID: req.User.ID},
	}

	err := c.db.Create(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] Create Category -1"
		log.Errorw(code, err)
		return err
	}
	return nil
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {
	var count int64
	err = c.db.Table("contents").Where("category_id  = ?", id).Count(&count).Error
	if err != nil {
		code = "[REPOSITORY] Delete Category -1"
		log.Errorw(code, err)
		return err
	}

	if count > 0 {
		return errors.New("category is used in content")
	}
	err = c.db.Where("id = ?", id).Delete(&model.Category{}).Error
	if err != nil {
		code = "[REPOSITORY] Delete Category -2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// EditCategory implements CategoryRepository.
func (c *categoryRepository) EditCategory(ctx context.Context, id int64, req entity.CategoryEntity) error {
	var count int64
	err = c.db.Table("categories").Where("slug = ?", req.Slug).Count(&count).Error
	if err != nil {
		code = "[REPOSITORY] Edit Category -2"
		log.Errorw(code, err)
		return err
	}

	slug := req.Slug
	if count == 0 {
		count = count + 1
		slug = fmt.Sprintf("%s-%d", req.Slug, count)
	}

	modelCategory := model.Category{
		Title: req.Title,
		Slug:  slug,
		User:  model.User{ID: req.User.ID},
	}

	err = c.db.Where("id = ?", id).Updates(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] Edit Category -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetCategory implements CategoryRepository.
func (c *categoryRepository) GetCategory(ctx context.Context) ([]entity.CategoryEntity, error) {
	var modelCategory []model.Category
	err := c.db.Order("created_at DESC").Preload("User").Find(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] Get Category -1"
		log.Errorw(code, err)
		return nil, err
	}

	if len(modelCategory) == 0 {
		code = "[REPOSITORY] Get Category -2"
		err = errors.New("category not found")
		log.Errorw(code, err)
		return nil, err
	}

	var res []entity.CategoryEntity
	for _, val := range modelCategory {
		res = append(res, entity.CategoryEntity{
			ID:    int(val.ID),
			Title: val.Title,
			Slug:  val.Slug,
			User: entity.UserEntity{
				ID:    int64(val.User.ID),
				Name:  val.User.Name,
				Email: val.User.Email,
			},
		})

	}
	return res, nil
}

// GetCategoryByID implements CategoryRepository.
func (c *categoryRepository) GetCategoryByID(ctx context.Context, id int64) (*entity.CategoryEntity, error) {
	var modelCategory model.Category
	err = c.db.Where("id = ?", id).Preload("User").First(&modelCategory).Error
	if err != nil {
		code = "[REPOSITORY] Get Category By ID -1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.CategoryEntity{
		ID:    int(modelCategory.ID),
		Title: modelCategory.Title,
		Slug:  modelCategory.Slug,
		User: entity.UserEntity{
			ID:    int64(modelCategory.User.ID),
			Name:  modelCategory.User.Name,
			Email: modelCategory.User.Email,
		},
	}, nil
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepository) UpdateCategory(ctx context.Context, id int64, req entity.CategoryEntity) (entity.CategoryEntity, error) {
	panic("unimplemented")
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
