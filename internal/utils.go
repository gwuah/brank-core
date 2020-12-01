package internal

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ConvertToUint64(num string) uint64 {
	id64, _ := strconv.ParseUint(num, 10, 64)
	return id64
}

func ConvertToInt(num string) int {
	id64, _ := strconv.ParseInt(num, 10, 64)
	return int(id64)
}
func Bool(b bool) *bool {
	temp := b
	return &temp
}
