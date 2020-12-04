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

		status, response, err := HandleMessagePost(req, r.eStore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
			})
		}
		c.JSON(status, response)
	})
}
