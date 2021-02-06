package routes

import (
	"brank/core"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterBrankRoutes(e *gin.Engine, s services.Services) {

	e.GET("/financial-institutions", func(c *gin.Context) {

		response := s.Brank.GetFinancialInstitutions()

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})

	e.POST("/message", func(c *gin.Context) {

		var req core.MessageRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.Brank.PublishMessageIntoKafka(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})
}
