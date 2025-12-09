package repository

import (
	"bwa-api/core/domain/entity"
	"bwa-api/core/domain/model"
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	UpdateContent(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

// CreateContent implements ContentRepository.
func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		CategoryID:  req.CategoryId,
		UserID:      req.UserID,
		Tags:        tags,
		Status:      req.Status,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Create Content -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// DeleteContent implements ContentRepository.
func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	err = c.db.Where("id = ?", id).Delete(&model.Content{}).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Delete Content -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

// GetContentByID implements ContentRepository.
func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content
	err := c.db.Where("id = ?", id).Preload(clause.Associations).First(&modelContent).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Get Content By ID -1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID:          modelContent.ID,
		Title:       modelContent.Title,
		Excerpt:     modelContent.Excerpt,
		Description: modelContent.Description,
		Image:       modelContent.Image,
		PublicId:    modelContent.PublicId,
		CategoryId:  modelContent.CategoryID,
		CreatedAt:   modelContent.CreatedAt,
		User: entity.UserEntity{
			ID:    modelContent.User.ID,
			Name:  modelContent.User.Name,
			Email: modelContent.User.Email,
		},
		Category: entity.CategoryEntity{
			ID:    int(modelContent.Category.ID),
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		Tags:   tags,
		Status: modelContent.Status,
	}

	return &resp, nil
}

// GetContents implements ContentRepository.
func (c *contentRepository) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, int64, int64, error) {
	var modelContents []model.Content
	var countData int64

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit
	status := ""
	if query.Status != "" {
		status = query.Status
	}

	sqlMain := c.db.Preload(clause.Associations).
		Where("title ilike ? OR excerpt ilike ? OR description ilike ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").Where("status LIKE  ?", "%"+status+"%")

	if query.CategoryId > 0 {
		sqlMain = sqlMain.Where("category_id = ?", query.CategoryId)
	}

	err = sqlMain.Model(&model.Content{}).Count(&countData).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Count Data -1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	err = sqlMain.Order(order).Offset(int(offset)).Limit(int(query.Limit)).Find(&modelContents).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Get Contents -1"
		log.Errorw(code, err)
		return nil, 0, 0, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Image:       val.Image,
			PublicId:    val.PublicId,
			CategoryId:  val.CategoryID,
			CreatedAt:   val.CreatedAt,
			User: entity.UserEntity{
				ID:    val.User.ID,
				Name:  val.User.Name,
				Email: val.User.Email,
			},
			Category: entity.CategoryEntity{
				ID:    int(val.Category.ID),
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			Tags:   tags,
			Status: val.Status,
		}

		resps = append(resps, resp)
	}
	return resps, countData, int64(totalPages), nil
}

// UpdateContent implements ContentRepository.
func (c *contentRepository) UpdateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		PublicId:    req.PublicId,
		CategoryID:  req.CategoryId,
		UserID:      req.User.ID,
		Tags:        tags,
		Status:      req.Status,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code = "[CONTENT REPOSITORY] Update Content -1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{db: db}

}
