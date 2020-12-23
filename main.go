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

	err = storage.RunMigrations(pg,
		&models.Transaction{},
		&models.Inquiry{},
		&models.Client{},
		&models.Bank{},
		&models.Account{},
		&models.Link{},
		&models.App{},
	)

	if config.RUN_SEEDS {
		log.Println("Running seeds")
		internal.RunSeeds(pg)
	}

	kvStore := storage.NewRedis(config)

	eventStore := internal.NewEventStore(config)
	repository := repository.NewRepo(pg)

	server := internal.NewHTTPServer(config)
	router := internal.NewRouter(server.Engine, eventStore, kvStore, repository, config)

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
