package repository

import "gorm.io/gorm"

type Repo struct {
	Transactions      *transactionLayer
	Clients           *clientLayer
	ClientApplication *clientApplicationLayer
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		Transactions:      newTransactionLayer(db),
		Clients:           newClientLayer(db),
		ClientApplication: newClientApplicationLayer(db),
	}
}
