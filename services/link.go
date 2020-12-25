package services

import (
	"brank/core"
	"brank/repository"
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type linkLayer struct {
	repo   repository.Repo
	config *core.Config
}

func newLinkLayer(r repository.Repo, c *core.Config) *linkLayer {
	return &linkLayer{
		repo:   r,
		config: c,
	}
}

func (l *linkLayer) ExchageContractCode(req core.ExchangeContractCode) core.BrankResponse {
	link, err := l.repo.Link.FindByCode(req.Code)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusUnauthorized,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "unauthorized request",
			},
		}
	}

	if err != nil {
		log.Println("link-FindByCode failed. err", err)
		return core.BrankResponse{
			Error: true,
			Code:  http.StatusBadRequest,
			Meta: core.BrankMeta{
				Data:    nil,
				Message: "account already exists. Wanna login?",
			},
		}
	}

	return core.BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: core.BrankMeta{
			Data: map[string]interface{}{
				"link": link,
			},
			Message: "exchange was successful",
		},
	}
}
