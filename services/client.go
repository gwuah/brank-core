package services

import (
	"brank/core"
	"brank/core/auth"
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
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusBadRequest,
			Meta: core.BrankMeta{
				Data:    client,
				Message: "Account already exists. Wanna login?",
			},
		}
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("findbyemail failed. err", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	passwordHash, err := auth.Hash(req.Password)
	if err != nil {
		log.Println("failed to hash password", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}
	client.FirstName = req.FirstName
	client.LastName = req.LastName
	client.CompanyName = req.CompanyName
	client.Password = passwordHash

	err = c.repo.Clients.Create(client)
	if err != nil {
		log.Println("create client failed. err", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	token, err := auth.GenerateClientAuthToken(client.ID, c.config.JWT_SIGNING_KEY)
	if err != nil {
		log.Println("failed to generate token", err)
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
				"token":  token,
				"client": client,
			},
			Message: "client created successfully",
		},
	}
}
