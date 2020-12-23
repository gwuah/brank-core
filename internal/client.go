package internal

import (
	"brank/internal/repository"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func CreateClient(req CreateClientRequest, e repository.Repo, c *Config) BrankResponse {
	client, err := e.Clients.FindByEmail(req.Email)
	if err == nil {
		return BrankResponse{
			Error: true,
			Code:  http.StatusBadRequest,
			Meta: BrankMeta{
				Data:    client,
				Message: "Account already exists. Wanna login?",
			},
		}
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("findbyemail failed. err", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	passwordHash, err := Hash(req.Password)
	if err != nil {
		log.Println("failed to hash password", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}
	client.FirstName = req.FirstName
	client.LastName = req.LastName
	client.CompanyName = req.CompanyName
	client.Password = passwordHash

	err = e.Clients.Create(client)
	if err != nil {
		log.Println("create client failed. err", err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusInternalServerError,
			Meta: BrankMeta{
				Data:    nil,
				Message: "request failed",
			},
		}
	}

	token, err := generateToken(fmt.Sprintf("%v", client.ID), c.JWT_SIGNING_KEY)
	if err != nil {
		log.Println("failed to generate token", err)
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
			Data:    client,
			Token:   token,
			Message: "client created successfully",
		},
	}
}
