package main

import (
	"brank/internal"
	"brank/internal/models"
	"brank/internal/storage"

	"log"
)

func main() {
	config := internal.NewConfig()

	pg, err := storage.NewPostgres(config)
	if err != nil {
		log.Fatal("postgres conn failed", err)
	}

	err = storage.RunMigrations(pg,
		&models.Client{},
	)

	pg.Create(&models.Client{Verified: internal.Bool(false)})

	redis := storage.NewRedis(config)

	if _, err := redis.Set("name", "manfo", 0).Result(); err != nil {
		log.Fatal("redis test insert error", err)
	}

	log.Println("shutting down")
}
