package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/utils"
	"brank/repository"
	"errors"
	"log"
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
		log.Println("failed to hash password", err)
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	client.FirstName = req.FirstName
	client.LastName = req.LastName
	client.CompanyName = req.CompanyName
	client.Password = passwordHash

	if err = c.repo.Clients.Create(client); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)

	}

	token, err := auth.GenerateClientAuthToken(client.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"token":  token,
		"client": client,
	}, utils.String("client created successfully"))

}
