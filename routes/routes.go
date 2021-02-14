package routes

import (
	"brank/core"
	"brank/core/auth"
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

func New(engine *gin.Engine, eStore mq.MQ, kvStore *redis.Client, repo repository.Repo, queue *queue.Que, config *core.Config, services services.Services) *router {
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
	RegisterBrankRoutes(r.engine, r.services)

	// App Links
	RegisterAppLinkRoutes(r.engine.Group("/links"), r.services)

	RegisterClientRoutes(r.engine.Group("/clients"), r.services)

	// Apps
	RegisterApplicationRoutes(r.engine.Group("/applications"), r.services)

	// Transactions
	RegisterTransactionsRoutes(r.engine.Group("/transactions", auth.AuthorizeProductRequest(r.config)), r.services)
}
