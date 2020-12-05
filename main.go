package main

import (
	"brank/internal"
	"brank/internal/models"
	"brank/internal/storage"
	"fmt"

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

	kvStore := storage.NewRedis(config)

	eventStore := internal.NewEventStore(config)

	server := internal.NewHTTPServer(config)
	router := internal.NewRouter(server.Engine, eventStore, kvStore)

	go func() {
		for {
			select {
			case msg := <-eventStore.Subscribe([]string{internal.GenerateTopic("validate_login")}):
				fmt.Println(string(msg), "from consumer")
			}
		}
	}()

	router.RegisterRoutes()
	server.Start()
}
