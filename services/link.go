package services

import (
	"brank/core"
	"brank/core/models"
	"brank/core/utils"
	"brank/integrations"
	"brank/repository"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type linkLayer struct {
	repo         repository.Repo
	config       *core.Config
	integrations integrations.Integrations
}

func newLinkLayer(r repository.Repo, c *core.Config, i integrations.Integrations) *linkLayer {
	return &linkLayer{
		repo:         r,
		config:       c,
		integrations: i,
	}
}

func (l *linkLayer) ExchageContractCode(req core.ExchangeContractCode) core.BrankResponse {
	link, err := l.repo.Link.FindByCode(req.Code)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Error(nil, utils.String("un-authorized"), http.StatusInternalServerError)
	}

	if err != nil {
		return utils.Error(nil, utils.String("account already exists. Wanna login?"), http.StatusBadRequest)
	}

	return utils.Success(&map[string]interface{}{
		"link": link,
	}, utils.String("exchange successful"))
}

func (l *linkLayer) LinkAccount(req core.LinkAccountRequest) core.BrankResponse {
	bank, err := l.repo.Bank.FindById(req.BankID)
	if err != nil {
		return utils.Error(nil, nil, http.StatusInternalServerError)
	}

	if bank.Code == models.FidelityBank {

	}

	return utils.Success(&map[string]interface{}{
		"bank": bank,
	}, utils.String("client created successfully"))
}
