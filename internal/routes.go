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
	config  *Config
}

func NewRouter(engine *gin.Engine, eStore EventStore, kvStore *redis.Client, repo repository.Repo, config *Config) *router {
	return &router{
		engine:  engine,
		eStore:  eStore,
		kvStore: kvStore,
		repo:    repo,
		config:  config,
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

		response := PublishMessageIntoKafka(req, r.eStore)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response)
	})

	r.engine.POST("/clients", func(c *gin.Context) {
		var req CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := CreateClient(req, r.repo, r.config)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	r.engine.POST("/client-application", func(c *gin.Context) {
		var req CreateAppRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := CreateApp(req, r.repo, r.config)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	r.engine.POST("/exchange", func(c *gin.Context) {
		var req ExchangeContractCode
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := ExchageContractCode(req, r.repo, r.config, r.eStore)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	r.engine.POST("/link-account", func(c *gin.Context) {
		var req LinkAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := LinkAccount(req, r.repo, r.config, r.eStore)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	// r.engine.GET("/transactions/:customerId", func(c *gin.Context) {
	// 	page := "1"
	// 	if c.Query("page") != "" {
	// 		page = c.Query("page")
	// 	}

	// 	response := HandleGetTransactions(TransactionsRequest{
	// 		CustomerId: ConvertToInt(c.Param("customerId")),
	// 		Page:       ConvertToInt(page),
	// 	}, r.repo)

	// 	if response.Error {
	// 		c.JSON(response.Code, gin.H{
	// 			"message": response.Meta.Message,
	// 		})
	// 		return
	// 	}

	// 	c.JSON(response.Code, response.Meta)
	// })
}
