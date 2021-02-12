package services

import (
	"brank/core"
	"brank/core/mq"
	"brank/core/queue"
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

func New(r repository.Repo, c *core.Config, mq mq.MQ, kv *redis.Client, q *queue.Que, i integrations.Integrations) Services {
	return Services{
		Clients:           newClientLayer(r, c),
		ClientApplication: newClientApplicationLayer(r, c),
		Links:             newLinkLayer(r, c, kv, q, i),
		Brank:             newBrankLayer(r, c, mq, q),
		Transactions:      newTransactionLayer(r, c),
	}
}
