package internal

import (
	"brank/internal/models"
	"brank/internal/repository"
	"log"
	"net/http"
	"strings"
)

func GenerateAccessToken(req CreateAppRequest, e repository.Repo, c *Config) BrankResponse {
	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    nil,
			Message: "client application created successfully",
		},
	}
}

func CreateApp(req CreateAppRequest, e repository.Repo, c *Config) BrankResponse {
	app := models.App{
		Name:        strings.ToLower(req.Name),
		Logo:        req.Logo,
		CallbackUrl: req.CallbackUrl,
		ClientID:    req.ClientID,
		PublicKey:   NewPublicKey(strings.ToLower(req.Name)),
	}
	err := e.ClientApplication.Create(&app)
	if err != nil {
		log.Println("create client_application failed. err", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}
	accessToken, err := generateAppAccessToken(app.ID, req.ClientID, c.JWT_SIGNING_KEY)
	if err != nil {
		log.Println("failed to generate access token", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	app.AccessToken = accessToken
	err = e.ClientApplication.Update(&app)
	if err != nil {
		log.Println("update client_application failed. err", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data: map[string]interface{}{
				"app": app,
			},
			Message: "client application created successfully",
		},
	}
}
