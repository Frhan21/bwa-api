package seeder

import (
	"bwa-api/core/domain/model"
	"bwa-api/libs/conv"

	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	bytes, err := conv.HashPassword("password")
	if err != nil {
		log.Fatal().Err(err).Msg("Error seeding admin user")
		return
	}

	admin := model.User{
		Name:     "Admin",
		Email:    "admin@admin.com",
		Password: string(bytes),
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: admin.Email}).Error; err != nil {
		log.Fatal().Err(err).Msg("Error seeding admin user")
		return
	} else {
		log.Info().Msg("Admin user seeded successfully")
	}
}
