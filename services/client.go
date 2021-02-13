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
}

func newClientLayer(r repository.Repo, c *core.Config) *clientLayer {
	return &clientLayer{
		repo:   r,
		config: c,
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

	passwordHash, err := auth.Hash(req.Password)
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

	token, err := auth.GenerateClientAccessToken(client.ID, c.config.JWT_SIGNING_KEY)
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

	status := auth.VerifyHash(client.Password, req.Password)
	if !status {
		return utils.Error(err, utils.String("email/password invalid"), http.StatusUnauthorized)
	}

	token, err := auth.GenerateClientAccessToken(client.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"token": token,
	}, utils.String("client login successful"))

}
