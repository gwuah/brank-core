package repository

import "gorm.io/gorm"

type Repo struct {
	Transactions      *transactionLayer
	Clients           *clientLayer
	ClientApplication *clientApplicationLayer
	Link              *linkLayer
	Bank              *bankLayer
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		Transactions:      newTransactionLayer(db),
		Clients:           newClientLayer(db),
		ClientApplication: newClientApplicationLayer(db),
		Link:              newLinkLayer(db),
		Bank:              newBankLayer(db),
	}
}
