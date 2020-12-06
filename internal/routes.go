package internal

import (
	"brank/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type router struct {
	engine  *gin.Engine
	kvStore *redis.Client
	eStore  EventStore
	repo    repository.Repo
}

func NewRouter(engine *gin.Engine, eStore EventStore, kvStore *redis.Client, repo repository.Repo) *router {
	return &router{
		engine:  engine,
		eStore:  eStore,
		kvStore: kvStore,
		repo:    repo,
	}
}

func (r *router) RegisterRoutes() {
	r.engine.POST("/message", func(c *gin.Context) {

		var req MessageRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := HandleMessagePost(req, r.eStore)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response)
	})

	r.engine.GET("/transactions/:customerId", func(c *gin.Context) {
		page := "1"
		if c.Query("page") != "" {
			page = c.Query("page")
		}

		response := HandleGetTransactions(TransactionsRequest{
			CustomerId: ConvertToInt(c.Param("customerId")),
			Page:       ConvertToInt(page),
		}, r.repo)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})
}
