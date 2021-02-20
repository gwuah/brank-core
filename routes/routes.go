package routes

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/mq"
	"brank/services"

	"brank/core/queue"
	"brank/routes/products"

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
	auth     *auth.Auth
}

func New(engine *gin.Engine, eStore mq.MQ, kvStore *redis.Client, repo repository.Repo, queue *queue.Que, config *core.Config, services services.Services, auth *auth.Auth) *router {
	return &router{
		engine:   engine,
		eStore:   eStore,
		kvStore:  kvStore,
		repo:     repo,
		config:   config,
		queue:    queue,
		services: services,
		auth:     auth,
	}
}

func (r *router) RegisterRoutes() {
	linksGroup := r.engine.Group("/links")
	clientsGroup := r.engine.Group("/clients")
	appsGroup := r.engine.Group("/applications")
	transactionsGroup := r.engine.Group("/transactions")

	// connect-routes
	RegisterAppLinkRoutes(linksGroup, r.services)

	// product routes
	products.RegisterTransactionsRoutes(transactionsGroup, r.services, r.auth)
	products.RegisterAppLinkRoutes(linksGroup, r.services, r.auth)

	// internal routes
	RegisterBrankRoutes(r.engine, r.services)
	RegisterClientRoutes(clientsGroup, r.services, r.auth)
	RegisterApplicationRoutes(appsGroup, r.services, r.auth)

}
