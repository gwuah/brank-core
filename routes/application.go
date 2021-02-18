package routes

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/utils"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterApplicationRoutes(e *gin.RouterGroup, s services.Services) {

	e.POST("", func(c *gin.Context) {
		var req core.CreateAppRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}

		response := s.Application.CreateApp(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.GET("/pk/:public_key", func(c *gin.Context) {
		publicKey := c.Param("public_key")

		if publicKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed",
			})
			return
		}

		response := s.Application.GetByPublicKey(publicKey)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.GET("/id/:id", func(c *gin.Context) {
		id := c.Param("id")

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed",
			})
			return
		}

		response := s.Application.GetByID(utils.ConvertToInt(id))

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.GET("", auth.AuthorizeClientRequest(s.Config), func(c *gin.Context) {
		clientId := c.GetInt("client_id")
		response := s.Application.All(clientId)
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
