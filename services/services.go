package services

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/mq"
	"brank/core/queue"
	"brank/integrations"
	"brank/repository"

	"github.com/go-redis/redis"
)

type Services struct {
	Config       *core.Config
	Clients      *clientLayer
	AppLinks     *appLinkLayer
	Application  *applicationLayer
	Brank        *brankLayer
	Transactions *transactionsLayer
}

func New(r repository.Repo, c *core.Config, mq mq.MQ, kv *redis.Client, q *queue.Que, i integrations.Integrations, a *auth.Auth) Services {
	return Services{
		Config:       c,
		Clients:      newClientLayer(r, c, a),
		Application:  newApplicationLayer(r, c, a),
		AppLinks:     newAppLinkLayer(r, c, kv, q, i, a),
		Brank:        newBrankLayer(r, c, mq, q),
		Transactions: newTransactionLayer(r, c),
	}
}
