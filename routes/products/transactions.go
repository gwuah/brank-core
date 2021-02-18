package products

import (
	"brank/core"
	"brank/core/auth"
	"brank/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterTransactionsRoutes(e *gin.RouterGroup, s services.Services) {
	e.POST("", func(c *gin.Context) {
		var req core.TransactionsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}

		appLinkId, err := auth.AuthorizeProductRequest(s.Config, req.AuthParams)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		response := s.Transactions.GetTransactions(services.GetTransactionsParams{
			AppLinkID: appLinkId,
			Offset:    req.Offset,
			Limit:     req.Limit,
		})
		if response.Error {
			c.JSON(response.Code, gin.H{
				"message": response.Meta.Message,
			})
			return
		}

		c.JSON(response.Code, response.Meta)
	})

}
