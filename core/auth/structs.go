package auth

import "github.com/dgrijalva/jwt-go"

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
