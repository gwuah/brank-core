package main

import (
	"brank/integrations"

	"brank/core"
	"brank/core/auth"
	serv "brank/core/server"

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
		&models.Client{},
		&models.Bank{},
		&models.Account{},
		&models.Link{},
		&models.App{},
		&models.Customer{},
		&models.AppLink{},
	)
	if err != nil {
		log.Fatal("failed to run migrations. err", err)
	}

	if c.RUN_SEEDS {
		log.Println("Running seeds")
		models.RunSeeds(pg)
	}

	cache := storage.NewRedis(c)
	r := repository.New(pg)

	i := integrations.New()

	mq, err := mq.New(c)
	if err != nil {
		log.Fatal("failed to initialize messaging queue. err", err)
	}

	q, err := queue.New(c)
	if err != nil {
		log.Fatal("failed to initialize queue. err", err)
	}
	workers := q.RegisterJobs(
		[]queue.JobWorker{
			que_workers.NewFidelityWorker(i, r, q),
			que_workers.NewFidelityTransactionsWorker(r),
			que_workers.NewWebhookWorker(r, q),
		},
	)
	go workers.Start()

	a := auth.New(r)
	s := services.New(r, c, mq, cache, q, *i, a)
	server := serv.NewHTTPServer(c, a)
	router := routes.New(server.Engine, mq, cache, r, q, c, s, a)

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
