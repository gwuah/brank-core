package products

import (
	"brank/core"
	"brank/core/auth"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAppLinkRoutes(e *gin.RouterGroup, s services.Services, a *auth.Auth) {
	e.POST("/exchange", func(c *gin.Context) {
		var req core.ExchangeContractCode
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}

		_, err := a.ParseAppToken(s.Config, core.AuthParams{
			AppLinkToken: req.AccessToken,
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		response := s.AppLinks.ExchageContractCode(req)

		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)

	})

}
