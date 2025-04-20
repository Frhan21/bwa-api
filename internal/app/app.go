package app

import (
	"bwa-api/config"
	"bwa-api/core/service"
	"bwa-api/internal/adapter/cloudflare"
	"bwa-api/internal/adapter/handler"
	"bwa-api/internal/adapter/repository"
	"bwa-api/libs/auth"
	"bwa-api/libs/middleware"
	"bwa-api/libs/pagination"
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"os"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgress()
	if err != nil {
		log.Fatal().Msgf("Error connection to database %v", err)
		return
	}

	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatal().Msgf("Error create directory %v", err)
		return
	}

	// CloudFlare
	cfR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cfR2)
	r2Adapter := cloudflare.NewCloudFlareR2Adapter(s3Client, cfg)
	jwt := auth.NewJwt(cfg)
	middlewareAuth := middleware.NewMiddleware(cfg)

	_ = pagination.Pagination()

	// Repository
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Service
	authServe := service.NewAuthService(authRepo, cfg, jwt)
	categoryServe := service.NewCategoryService(categoryRepo)
	contentServe := service.NewContentService(contentRepo, cfg, r2Adapter)
	userServe := service.NewUserService(userRepo)

	// Handler
	authHandler := handler.NewAuthHandler(authServe)
	categoryHandler := handler.NewCategoryHandler(categoryServe)
	contentHandler := handler.NewContentHandler(contentServe)
	userHandler := handler.NewUserHandler(userServe)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		cfg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API Docs",
		}

		app.Use(swagger.New(cfg))
	}
	api := app.Group("/api")
	api.Post("/login", authHandler.Login)

	adminApp := api.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	// Category
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategory)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Get("/:id", categoryHandler.GetCategoryByID)
	categoryApp.Put("/:id", categoryHandler.EditCategory)
	categoryApp.Delete("/:id", categoryHandler.DeleteCategory)

	// Content
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Get("/:id", contentHandler.GetContentByID)
	contentApp.Put("/:id", contentHandler.UpdateContent)
	contentApp.Delete("/:id", contentHandler.DeleteContent)
	contentApp.Post("/upload", contentHandler.UploadImage)

	// User
	userApp := adminApp.Group("/users")
	userApp.Get("/profile", userHandler.GetUserById)
	userApp.Put("/update-password", userHandler.UpdatePassword)

	// FE Route
	feApp := api.Group("/fe")
	feApp.Get("/categories", categoryHandler.GetCategoryFE)
	feApp.Get("/contents", contentHandler.GetContentWithQuery)
	feApp.Get("/contents/:id", contentHandler.GetContentDetail)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal().Msgf("Error start server %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	<-quit
	log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
