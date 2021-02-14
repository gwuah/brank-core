package routes

import (
	"brank/core"
	"brank/core/utils"
	"brank/services"

	"github.com/gin-gonic/gin"
)

// func extractPaginationDataFromUrl(c *gin.Context) {
// 	var (
// 		offset = c.Query("offset")
// 		limit  = c.Query("limit")
// 	)
// }

func RegisterTransactionsRoutes(e *gin.RouterGroup, s services.Services) {
	e.GET("", func(c *gin.Context) {
		page := "1"
		if c.Query("offset") != "" {
			page = c.Query("page")
		}

		response := s.Transactions.GetTransactions(core.TransactionsRequest{
			CustomerId: utils.ConvertToInt(c.Param("accountId")),
			Page:       utils.ConvertToInt(page),
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
