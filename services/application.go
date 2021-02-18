package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/models"
	"brank/core/utils"
	"brank/repository"
	"net/http"
	"strings"
)

type applicationLayer struct {
	repo   repository.Repo
	config *core.Config
}

func newApplicationLayer(r repository.Repo, c *core.Config) *applicationLayer {
	return &applicationLayer{
		repo:   r,
		config: c,
	}
}

func (a *applicationLayer) CreateApp(req core.CreateAppRequest) core.BrankResponse {
	app := models.App{
		Name:        strings.ToLower(req.Name),
		Logo:        req.Logo,
		CallbackUrl: req.CallbackUrl,
		ClientID:    req.ClientID,
		PublicKey:   utils.NewPublicKey(strings.ToLower(req.Name)),
	}
	if err := a.repo.Application.Create(&app); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	accessToken, err := auth.GenerateAppAccessToken(app.ID, a.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	app.AccessToken = accessToken
	if err := a.repo.Application.Update(&app); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"app": app,
	}, nil)
}

func (a *applicationLayer) GetByPublicKey(key string) core.BrankResponse {
	app, err := a.repo.Application.FindByPublicKey(key)
	if err != nil {
		return utils.Error(err, utils.String("not found"), http.StatusNotFound)
	}
	return utils.Success(&map[string]interface{}{
		"app": app,
	}, nil)
}

func (a *applicationLayer) GetByID(id int) core.BrankResponse {
	app, err := a.repo.Application.FindByID(id)
	if err != nil {
		return utils.Error(err, utils.String("not found"), http.StatusNotFound)
	}
	return utils.Success(&map[string]interface{}{
		"app": app,
	}, nil)
}

func (a *applicationLayer) All(id int) core.BrankResponse {
	apps, err := a.repo.Application.All("client_id=?", id)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}
	return utils.Success(&map[string]interface{}{
		"apps": apps,
	}, nil)
}
