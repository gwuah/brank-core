package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	AppID    int `json:"app_id"`
	ClientID int `json:"client_id"`
	jwt.StandardClaims
}

func Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func VerifyHash(hash, value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)) == nil
}

func GenerateClientAuthToken(ClientID int, key string) (string, error) {
	const expirationHours = 24 * 90
	claims := CustomClaims{
		ClientID: ClientID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expirationHours).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}

func GenerateAppAccessToken(appId int, clientId int, key string) (string, error) {
	claims := CustomClaims{
		AppID:    appId,
		ClientID: clientId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}
