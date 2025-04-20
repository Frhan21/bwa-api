package cloudflare

import (
	"bwa-api/config"
	"bwa-api/core/domain/entity"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/log"
)

var code string
var err error

type CloudFlareR2Adapter interface {
	UploadImage(req *entity.FileUploadRequest) (string, error)
}

// UploadImage uploads an image to the Cloudflare R2 bucket.
func (c *cloudFlareR2Adapter) UploadImage(req *entity.FileUploadRequest) (string, error) {
	openedFile, err := os.Open(req.Path)
	if err != nil {
		code = "[CLOUDFLARE R2 ADAPTER] Upload Image -1"
		log.Errorw(code, err)
		return "", err
	}

	defer openedFile.Close()
	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.Bucket),
		Key:         aws.String(req.Name),
		Body:        openedFile,
		ContentType: aws.String("image/jpeg"),
	})

	if err != nil {
		code = "[CLOUDFLARE R2 ADAPTER] Upload Image -2"
		log.Errorw(code, err)
		return "", err
	}

	return fmt.Sprintf("%s/%s/%s", c.BaseUrl, c.Bucket, req.Name), nil
}

type cloudFlareR2Adapter struct {
	Client  *s3.Client
	Bucket  string
	BaseUrl string
}

func NewCloudFlareR2Adapter(client *s3.Client, cfg *config.Config) CloudFlareR2Adapter {
	clientBase := s3.NewFromConfig(cfg.LoadAwsConfig(), func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.CF.AccountId))
	})
	return &cloudFlareR2Adapter{
		Client:  clientBase,
		Bucket:  cfg.CF.Name,
		BaseUrl: cfg.CF.PublicUrl,
	}
}
