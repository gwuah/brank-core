package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type ClientClaim struct {
	ClientID int `json:"client_id"`
	jwt.StandardClaims
}

type AppClaim struct {
	AppID int `json:"app_id"`
	jwt.StandardClaims
}

type ExchangeClaim struct {
	AppLinkID int `json:"app_link_id"`
	jwt.StandardClaims
}

func Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func VerifyHash(hash, value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)) == nil
}

func GenerateClientAccessToken(ClientID int, key string) (string, error) {
	const expirationHours = 24
	claims := ClientClaim{
		ClientID: ClientID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expirationHours).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}

func GenerateAppAccessToken(appId int, key string) (string, error) {
	claims := AppClaim{
		AppID: appId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}

func GenerateExchangeAccessToken(appLinkID int, key string) (string, error) {
	const expirationHours = 24 * 30
	claims := ExchangeClaim{
		AppLinkID: appLinkID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * expirationHours).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}
