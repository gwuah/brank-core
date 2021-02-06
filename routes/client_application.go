package routes

import (
	"brank/core"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterClientApplicationRoutes(e *gin.RouterGroup, s services.Services) {

	e.POST("", func(c *gin.Context) {
		var req core.CreateAppRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.ClientApplication.CreateApp(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.GET("/:public_key", func(c *gin.Context) {
		publicKey := c.Param("public_key")

		response := s.ClientApplication.GetByPublicKey(publicKey)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
