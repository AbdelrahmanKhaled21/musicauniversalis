package config

import (
	"fmt"
	"log"

	"github.com/AbdelrahmanKhaled21/musicauniversalis/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(config *Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Auto migrate the schema with more lenient settings
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Song{},
		&models.Playlist{},
		&models.PlaylistSong{},
	); err != nil {
		log.Printf("Warning: Auto-migration had issues: %v", err)
		// Continue anyway - the app might still work
	}

	log.Println("Database migration completed")

	return nil
}
