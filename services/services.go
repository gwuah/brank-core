package services

import (
	"brank/core"
	"brank/core/mq"
	"brank/integrations"
	"brank/repository"

	"github.com/go-redis/redis"
)

type Services struct {
	Clients           *clientLayer
	Links             *linkLayer
	ClientApplication *clientApplicationLayer
	Brank             *brankLayer
	Transactions      *transactionsLayer
}

func NewService(r repository.Repo, c *core.Config, mq mq.MQ, kv *redis.Client, i integrations.Integrations) Services {
	return Services{
		Clients:           newClientLayer(r, c),
		ClientApplication: newClientApplicationLayer(r, c),
		Links:             newLinkLayer(r, c, kv, i),
		Brank:             newBrankLayer(r, c, mq),
		Transactions:      newTransactionLayer(r, c),
	}
}
