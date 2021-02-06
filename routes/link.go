package routes

import (
	"brank/core"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterLinkRoutes(e *gin.RouterGroup, s services.Services) {

	e.POST("", func(c *gin.Context) {
		var req core.LinkAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.Links.LinkAccount(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

	e.POST("/otp", func(c *gin.Context) {
		var req core.VerifyOTPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "bad request",
			})
			return
		}

		response := s.Links.VerifyOTP(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

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
