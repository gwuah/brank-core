package main

import (
	"brank/integrations"

	"brank/core"
	"brank/core/models"
	"brank/core/mq"
	"brank/core/queue"
	"brank/core/storage"
	"brank/core/utils"
	que_workers "brank/core/workers"

	"brank/repository"
	"brank/routes"
	"brank/services"
	"fmt"

	"log"
)

func main() {
	c := core.NewConfig()
	pg, err := storage.NewPostgres(c)
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
		&models.Customer{},
	)
	if err != nil {
		log.Fatal("failed to run migrations. err", err)
	}

	if c.RUN_SEEDS {
		log.Println("Running seeds")
		models.RunSeeds(pg)
	}

	cache := storage.NewRedis(c)
	r := repository.NewRepo(pg)

	integrations := integrations.NewBankIntegrations()

	mq, err := mq.NewMQ(c)
	if err != nil {
		log.Fatal("failed to initialize messaging queue. err", err)
	}

	q, err := queue.NewQue(c)
	if err != nil {
		log.Fatal("failed to initialize queue. err", err)
	}
	workers := q.RegisterJobs(
		[]queue.JobWorker{
			que_workers.NewFidelityWorker(),
		},
	)
	go workers.Start()

	s := services.NewService(r, c, mq, *integrations)
	server := core.NewHTTPServer(c)
	router := routes.NewRouter(server.Engine, mq, cache, r, q, c, s)

	go func() {
		stream := mq.Subscribe([]string{utils.GenerateTopic("validate_login")})
		for {
			select {
			case msg := <-stream:
				fmt.Println(string(msg), "from consumer")
			}
		}
	}()

	router.RegisterRoutes()
	server.Start(func() {
		workers.Shutdown()
		q.Close()
	})

}
