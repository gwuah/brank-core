package services

import (
	"brank/core"
	"brank/core/mq"
	"brank/repository"
)

type Services struct {
	Clients           *clientLayer
	Links             *linkLayer
	ClientApplication *clientApplicationLayer
	Brank             *brankLayer
	Transactions      *transactionsLayer
}

func NewService(r repository.Repo, c *core.Config, mq mq.MQ) Services {
	return Services{
		Clients:           newClientLayer(r, c),
		ClientApplication: newClientApplicationLayer(r, c),
		Links:             newLinkLayer(r, c),
		Brank:             newBrankLayer(r, c, mq),
		Transactions:      newTransactionLayer(r, c),
	}
}
