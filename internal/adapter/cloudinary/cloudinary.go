package cloudinary

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2/log"
)

type CloudinaryAdapter interface {
	UploadImage(ctx context.Context, req *entity.FileUploadRequest) (*entity.FileUploadResponse, error)
	DeleteImage(ctx context.Context, publicId string) error
}

type cloudinaryAdapter struct {
	cld    *cloudinary.Cloudinary
	folder string
}

func (c *cloudinaryAdapter) UploadImage(ctx context.Context, req *entity.FileUploadRequest) (*entity.FileUploadResponse, error) {
	uploadResult, err := c.cld.Upload.Upload(ctx, req.Path, uploader.UploadParams{
		PublicID: req.Name,
		Folder:   c.folder,
	})

	if err != nil {
		log.Errorw("[CLOUDINARY-1] Upload Image", err)
		return nil, err
	}
	return &entity.FileUploadResponse{
		Url:      uploadResult.SecureURL,
		PublicId: uploadResult.PublicID,
	}, nil
}

func (c *cloudinaryAdapter) DeleteImage(ctx context.Context, publicId string) error {
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicId,
	})
	if err != nil {
		log.Errorw("[CLOUDINARY ADAPTER] Delete Image", err)
		return err
	}
	return nil
}

func NewCloudinaryAdapter(cfg *config.Config) CloudinaryAdapter {
	cld, err := cloudinary.NewFromParams(cfg.CD.Cloudname, cfg.CD.Apikey, cfg.CD.ApiSecret)
	if err != nil {
		log.Fatalw("[CLOUDINARY-2] Init Cloudinary", err)
	}

	return &cloudinaryAdapter{
		cld:    cld,
		folder: cfg.CD.UploadFile,
	}
}
