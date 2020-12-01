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

	_ = storage.NewRedis(config)

	server := internal.NewHTTPServer(config)
	server.Start()
}
