package clients

import (
	"fmt"

	"github.com/fazilnbr/go-clean-architecture/internal/app/users"
	"github.com/fazilnbr/go-clean-architecture/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(config utils.Postgres) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.Host, config.User, config.Password, config.Database, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Verify the connection is valid
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting db instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %w", err)
	}

	// create user table according to the struct and its fields
	db.AutoMigrate(users.User{})

	return db, nil
}
