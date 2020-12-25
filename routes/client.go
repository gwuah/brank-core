package routes

import (
	"brank/core"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(e *gin.Engine, s services.Services) {

	e.POST("/clients", func(c *gin.Context) {
		var req core.CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.Clients.CreateClient(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
