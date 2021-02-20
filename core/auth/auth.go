package auth

import (
	"brank/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	r repository.Repo
}

func New(r repository.Repo) *Auth {
	return &Auth{r: r}
}

func (a Auth) Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 10)
	return string(bytes), err
}

func (a Auth) VerifyHash(hash, value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(value)) == nil
}

func (a Auth) GenerateClientAccessToken(ClientID int, key string) (string, error) {
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

func (a Auth) GenerateAppAccessToken(appId int, key string) (string, error) {
	claims := AppClaim{
		AppID: appId,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}

func (a Auth) GenerateExchangeAccessToken(appLinkID int, key string) (string, error) {
	const expirationHours = 24 * 30
	claims := ExchangeClaim{
		AppLinkID: appLinkID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(key))
}
