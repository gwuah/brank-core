package routes

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/utils"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(e *gin.RouterGroup, s services.Services, a *auth.Auth) {

	e.GET("/:id", a.AuthorizeClientRequest(s.Config), func(c *gin.Context) {
		id := c.Param("id")

		response := s.Clients.Get(utils.ConvertToInt(id))

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.POST("", func(c *gin.Context) {
		var req core.CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
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

	e.POST("/login", func(c *gin.Context) {
		var req core.LoginClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}

		response := s.Clients.Login(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
