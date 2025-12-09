package config

import (
	"bwa-api/database/seeder"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgress struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgress() (*Postgress, error) {

	sslMode := cfg.PsqlDB.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.PsqlDB.User, cfg.PsqlDB.Password, cfg.PsqlDB.Host, cfg.PsqlDB.Port, cfg.PsqlDB.DBName, sslMode)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{}) // Use gorm.io/driver/postgres
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	seeder.SeedRoles(db)

	sqlDB.SetMaxOpenConns(cfg.PsqlDB.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.PsqlDB.DBMaxIdle)

	return &Postgress{DB: db}, nil

}
