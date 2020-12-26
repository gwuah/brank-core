package routes

import (
	"brank/core"
	"brank/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterLinkRoutes(e *gin.Engine, s services.Services) {
	log.Println("nothing to commit, working tree clean")
	e.POST("/exchange", func(c *gin.Context) {
		var req core.ExchangeContractCode
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.Links.ExchageContractCode(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
