package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/utils"
	"brank/repository"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type clientLayer struct {
	repo   repository.Repo
	config *core.Config
	auth   *auth.Auth
}

func newClientLayer(r repository.Repo, c *core.Config, a *auth.Auth) *clientLayer {
	return &clientLayer{
		repo:   r,
		config: c,
		auth:   a,
	}
}

func (c *clientLayer) CreateClient(req core.CreateClientRequest) core.BrankResponse {
	client, err := c.repo.Clients.FindByEmail(req.Email)
	if err == nil {
		return utils.Error(err, utils.String("account already exists"), http.StatusBadRequest)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	passwordHash, err := c.auth.Hash(req.Password)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	client.FirstName = req.FirstName
	client.LastName = req.LastName
	client.CompanyName = req.CompanyName
	client.Password = passwordHash
	client.Verified = utils.Bool(false)

	if err = c.repo.Clients.Create(client); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	token, err := c.auth.GenerateClientAccessToken(client.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"token":  token,
		"client": client,
	}, utils.String("client created successfully"))

}

func (c *clientLayer) Login(req core.LoginClientRequest) core.BrankResponse {
	client, err := c.repo.Clients.FindByEmail(req.Email)
	if err != nil {
		return utils.Error(err, utils.String("email/password invalid"), http.StatusUnauthorized)
	}

	status := c.auth.VerifyHash(client.Password, req.Password)
	if !status {
		return utils.Error(err, utils.String("email/password invalid"), http.StatusUnauthorized)
	}

	token, err := c.auth.GenerateClientAccessToken(client.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"token":  token,
		"client": client,
	}, nil)

}

func (c *clientLayer) Get(id int) core.BrankResponse {
	client, err := c.repo.Clients.FindByID(id)
	if err != nil {
		return utils.Error(err, utils.String("not found"), http.StatusNotFound)
	}
	return utils.Success(&map[string]interface{}{
		"client": client,
	}, nil)
}
