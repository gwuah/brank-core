package main

import (
	"brank/internal"
	"brank/internal/models"
	"brank/internal/repository"
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

	err = storage.SetupJoinTable(pg, &models.Customer{}, "Banks", &models.CustomerBank{})
	err = storage.SetupJoinTable(pg, &models.Bank{}, "Customers", &models.CustomerBank{})

	err = storage.RunMigrations(pg,
		&models.Transaction{},
		&models.Customer{},
		&models.Inquiry{},
		&models.Client{},
		&models.Bank{},
	)

	kvStore := storage.NewRedis(config)

	eventStore := internal.NewEventStore(config)
	repository := repository.NewRepo(pg)

	server := internal.NewHTTPServer(config)
	router := internal.NewRouter(server.Engine, eventStore, kvStore, repository)

	go func() {
		stream := eventStore.Subscribe([]string{internal.GenerateTopic("validate_login")})
		for {
			select {
			case msg := <-stream:
				fmt.Println(string(msg), "from consumer")
			}
		}
	}()

	router.RegisterRoutes()
	server.Start()
}
