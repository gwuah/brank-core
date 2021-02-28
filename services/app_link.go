package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/models"
	"brank/core/queue"
	"brank/core/utils"
	worker "brank/core/workers"
	"brank/integrations"
	"brank/repository"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type appLinkLayer struct {
	cache        *redis.Client
	repo         repository.Repo
	config       *core.Config
	integrations integrations.Integrations
	q            *queue.Que
	auth         *auth.Auth
}

func newAppLinkLayer(r repository.Repo, c *core.Config, kv *redis.Client, q *queue.Que, i integrations.Integrations, a *auth.Auth) *appLinkLayer {
	return &appLinkLayer{
		repo:         r,
		config:       c,
		integrations: i,
		cache:        kv,
		q:            q,
		auth:         a,
	}
}

func (l *appLinkLayer) ExchageContractCode(req core.ExchangeContractCode) core.BrankResponse {
	appLink, err := l.repo.AppLink.FindByCode(req.Code)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Error(nil, nil, http.StatusUnauthorized)
	}

	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	accessToken, err := l.auth.GenerateExchangeAccessToken(appLink.ID, l.config.JWT_SIGNING_KEY)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	appLink.AccessToken = accessToken
	appLink.State = models.Claimed

	if err = l.repo.AppLink.Update(appLink); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"access_token": accessToken,
	}, nil)
}

func (l *appLinkLayer) LinkAccount(req core.LinkAccountRequest) core.BrankResponse {
	bank, err := l.repo.Bank.FindById(req.BankID)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	switch bank.Code {
	case models.FidelityBank:
		response, err := l.integrations.Fidelity.VerifyLogin(req.Username, req.Password)
		if err != nil {
			return utils.Error(err, utils.String("login failed"), http.StatusUnauthorized)
		}

		app, err := l.repo.Application.FindByPublicKey(req.PublicKey)
		if err != nil {
			return utils.Error(err, utils.String("public_key is invalid"), http.StatusUnauthorized)
		}

		link := models.Link{
			BankID:   bank.ID,
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

		appLink := models.AppLink{
			AppID:  app.ID,
			LinkID: link.ID,
			State:  models.Unclaimed,
			Code:   utils.GenerateExchangeCode(),
		}

		if err := l.repo.AppLink.Create(&appLink); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		sessionId := utils.GenerateUUID()

		if err := l.cache.Set(sessionId, appLink.ID, 0).Err(); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		return utils.Success(&map[string]interface{}{
			"session_id":   sessionId,
			"requires_otp": bank.RequiresOtp,
		}, utils.String("link sucessful"))

	}

	return utils.Error(nil, utils.String("no integration available for bank"), http.StatusBadRequest)

}

func (l *appLinkLayer) VerifyOTP(req core.VerifyOTPRequest) core.BrankResponse {
	val, err := l.cache.Get(fmt.Sprint(req.SessionID)).Result()
	if err == redis.Nil {
		return utils.Error(err, utils.String("session_id has either expired or is invalid"), http.StatusBadRequest)
	}

	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	appLink, err := l.repo.AppLink.FindById(utils.ConvertToInt(val))
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	link, err := l.repo.Link.FindById(appLink.LinkID)
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
		response, err := l.integrations.Fidelity.VerifyOtp(req.Otp)
		if err != nil {
			return utils.Error(nil, utils.String("otp verification failed"), http.StatusUnauthorized)
		}

		meta.Fidelity.Otp = *response

		if err := link.CommitMeta(meta); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		if err := l.repo.Link.Update(link); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		if err := l.q.QueueJob(worker.FidelityJob, worker.CreateFidelityJob(appLink.ID)); err != nil {
			return utils.Error(err, nil, http.StatusInternalServerError)
		}

		return utils.Success(&map[string]interface{}{
			"code": appLink.Code,
		}, utils.String("link complete"))

	}

	return utils.Error(nil, utils.String("no integration available for bank"), http.StatusBadRequest)

}
