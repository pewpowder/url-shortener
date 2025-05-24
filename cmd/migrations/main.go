package main

import (
	"log"

	"github.com/pewpowder/url-shortener/internal/config"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/pewpowder/url-shortener/internal/resources"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, _ := config.MustLoad(resources.AppConfigFS)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.DB.DSN,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Can't connect to DB: %s", err)
	}

	if err := db.AutoMigrate(&entity.Link{}, &entity.LinkStats{}, &entity.Visit{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Default().Println("Migrations executed successfully")
}
