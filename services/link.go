package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/models"
	"brank/core/utils"
	"brank/integrations"
	"brank/repository"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type linkLayer struct {
	cache        *redis.Client
	repo         repository.Repo
	config       *core.Config
	integrations integrations.Integrations
}

func newLinkLayer(r repository.Repo, c *core.Config, kv *redis.Client, i integrations.Integrations) *linkLayer {
	return &linkLayer{
		repo:         r,
		config:       c,
		integrations: i,
		cache:        kv,
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
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	switch bank.Code {
	case models.FidelityBank:
		status, response, err := l.integrations.Fidelity.VerifyLogin(req.Username, req.Password)
		if err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		if status {
			app, err := l.repo.ClientApplication.FindByPublicKey(req.PublicKey)
			if err != nil {
				return utils.Error(err, nil, http.StatusInternalServerError)
			}

			link := models.Link{
				Code:     auth.GenerateExchangeCode(),
				BankID:   bank.ID,
				AppID:    app.ID,
				Username: req.Username,
				Password: req.Password,
			}

			meta, err := link.GetMeta()
			if err != nil {
				return utils.Error(err, nil, http.StatusInternalServerError)
			}

			meta.Fidelity = models.Fidelity{
				Init: *response,
			}

			if err := link.CommitMeta(meta); err != nil {
				return utils.Error(err, nil, http.StatusInternalServerError)
			}

			if err := l.repo.Link.Create(&link); err != nil {
				return utils.Error(err, nil, http.StatusInternalServerError)
			}

			sessionId := utils.GenerateUUID()

			if err := l.cache.Set(sessionId, link.ID, 0).Err(); err != nil {
				return utils.Error(err, nil, http.StatusInternalServerError)
			}

			return utils.Success(&map[string]interface{}{
				"session_id":   sessionId,
				"requires_otp": bank.RequiresOtp,
			}, utils.String("link sucessful"))
		}

		return utils.Error(nil, nil, http.StatusUnauthorized)
	}

	return utils.Error(nil, utils.String("no integration available for bank"), http.StatusBadRequest)

}

func (l *linkLayer) VerifyOTP(req core.VerifyOTPRequest) core.BrankResponse {
	val, err := l.cache.Get(fmt.Sprint(req.SessionID)).Result()
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	link, err := l.repo.Link.FindById(utils.ConvertToInt(val))
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	bank, err := l.repo.Bank.FindById(link.BankID)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	switch bank.Code {
	case models.FidelityBank:

		meta, err := link.GetMeta()
		if err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		l.integrations.Fidelity.SetBearerToken(meta.Fidelity.Init.Token)
		status, response, err := l.integrations.Fidelity.VerifyOtp(req.Otp)
		if err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		meta.Fidelity.Otp = *response

		if err := link.CommitMeta(meta); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		if status {
			return utils.Success(&map[string]interface{}{
				"code": link.Code,
			}, utils.String("otp verification successful"))
		}

		return utils.Error(nil, nil, http.StatusUnauthorized)
	}

	return utils.Error(nil, utils.String("no integration available for bank"), http.StatusBadRequest)

}
