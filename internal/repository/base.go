package repository

import "gorm.io/gorm"

type Repo struct {
	Transactions *transactionLayer
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		Transactions: newTransactionLayer(db),
	}
}
