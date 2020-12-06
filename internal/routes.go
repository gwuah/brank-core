package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type router struct {
	engine  *gin.Engine
	kvStore *redis.Client
	eStore  EventStore
}

func NewRouter(engine *gin.Engine, eStore EventStore, kvStore *redis.Client) *router {
	return &router{
		engine:  engine,
		eStore:  eStore,
		kvStore: kvStore,
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

	r.engine.GET("/transactions/:userId", func(c *gin.Context) {
		page := "1"
		if c.Query("page") != "" {
			page = c.Query("page")
		}

		response := HandleGetTransactions(TransactionsRequest{
			UserId: c.Param("userId"),
			Page:   ConvertToInt(page),
		})

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})
}
