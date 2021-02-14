package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/models"
	"brank/core/utils"
	"brank/repository"
	"log"
	"net/http"
	"strings"
)

type clientApplicationLayer struct {
	repo   repository.Repo
	config *core.Config
}

func newClientApplicationLayer(r repository.Repo, c *core.Config) *clientApplicationLayer {
	return &clientApplicationLayer{
		repo:   r,
		config: c,
	}
}

func (c *clientApplicationLayer) CreateApp(req core.CreateAppRequest) core.BrankResponse {
	app := models.App{
		Name:        strings.ToLower(req.Name),
		Logo:        req.Logo,
		CallbackUrl: req.CallbackUrl,
		ClientID:    req.ClientID,
		PublicKey:   utils.NewPublicKey(strings.ToLower(req.Name)),
	}
	err := c.repo.ClientApplication.Create(&app)
	if err != nil {
		log.Println("create client_application failed. err", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}
	accessToken, err := auth.GenerateAppAccessToken(app.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		log.Println("failed to generate access token", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	app.AccessToken = accessToken
	err = c.repo.ClientApplication.Update(&app)
	if err != nil {
		log.Println("update client_application failed. err", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	return core.BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: core.BrankMeta{
			Data: map[string]interface{}{
				"app": app,
			},
			Message: "client application created successfully",
		},
	}
}

func (c *clientApplicationLayer) GetByPublicKey(key string) core.BrankResponse {
	app, err := c.repo.ClientApplication.FindByPublicKey(key)
	if err != nil {
		return utils.Error(err, utils.String("not found"), http.StatusNotFound)
	}
	return utils.Success(&map[string]interface{}{
		"app": app,
	}, nil)
}
