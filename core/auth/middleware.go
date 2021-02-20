package auth

import (
	"brank/core"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func extractToken(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], "bearer") {
		return "", false
	}
	return authHeaderParts[1], true
}

func (a Auth) CORS() gin.HandlerFunc {
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

func (a Auth) ExtractTokenFromAuthHeader(cfg *core.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractToken(c.Request.Header.Get("Authorization"))
		if ok {
			c.Set("token", token)
		}
		c.Next()
	}
}

func (a Auth) ParseAppToken(cfg *core.Config, req core.AuthParams) (int, error) {
	var claim AppClaim
	_, err := jwt.ParseWithClaims(req.AppLinkToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT_SIGNING_KEY), nil
	})
	if err != nil {
		return 0, err
	}
	return claim.AppID, nil
}

func (a Auth) AuthorizeProductRequest(cfg *core.Config, req core.AuthParams) (int, error) {
	claim := ExchangeClaim{}
	_, err := jwt.ParseWithClaims(req.AppLinkToken, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT_SIGNING_KEY), nil
	})
	if err != nil {
		return 0, err
	}
	return claim.AppLinkID, nil
}

func (a Auth) AuthorizeClientRequest(cfg *core.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		tokenString, _ := token.(string)
		claim := ClientClaim{}
		_, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT_SIGNING_KEY), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		c.Set("client_id", claim.ClientID)
	}
}
