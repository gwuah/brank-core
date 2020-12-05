package internal

import (
	"fmt"
	"net/http"
	"time"

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
		time1 := time.Now()

		var req MessageRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		fmt.Println("Time since after parsing", time.Since(time1))
		status, response, err := HandleMessagePost(req, r.eStore)
		fmt.Println("Time since after service runs", time.Since(time1))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
			})
			return
		}

		fmt.Println("Time since before response runs", time.Since(time1))

		c.JSON(status, response)
	})
}
