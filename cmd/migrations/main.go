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
		log.Fatalf("can't connect to DB for migrations: %s", err)
	}

	if err := db.AutoMigrate(&entity.Link{}, &entity.LinkStats{}, &entity.Visit{}, &entity.Tag{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Default().Println("migrations executed successfully")
}
