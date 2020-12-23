package internal

import (
	"brank/internal/models"
	"brank/internal/repository"
	"log"
	"net/http"
)

func CreateApp(req CreateAppRequest, e repository.Repo, c *Config) BrankResponse {
	app := models.App{
		Name:        req.Name,
		Logo:        req.Logo,
		CallbackUrl: req.CallbackUrl,
		ClientID:    req.ClientID,
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

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    app,
			Message: "client application created successfully",
		},
	}
}
