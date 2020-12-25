package routes

import (
	"brank/core"
	"brank/core/mq"
	"brank/services"

	"brank/core/queue"
	"brank/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type router struct {
	engine   *gin.Engine
	kvStore  *redis.Client
	eStore   mq.MQ
	repo     repository.Repo
	queue    *queue.Que
	config   *core.Config
	services services.Services
}

func NewRouter(engine *gin.Engine, eStore mq.MQ, kvStore *redis.Client, repo repository.Repo, queue *queue.Que, config *core.Config, services services.Services) *router {
	return &router{
		engine:   engine,
		eStore:   eStore,
		kvStore:  kvStore,
		repo:     repo,
		config:   config,
		queue:    queue,
		services: services,
	}
}

func (r *router) RegisterRoutes() {
	RegisterLinkRoutes(r.engine, r.services)
	RegisterClientRoutes(r.engine, r.services)
	RegisterClientApplicationRoutes(r.engine, r.services)
	RegisterBrankRoutes(r.engine, r.services)
	RegisterTransactionsRoutes(r.engine, r.services)
}