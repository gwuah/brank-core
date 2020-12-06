package repository

import "gorm.io/gorm"

type Repo struct {
	Customers    *customerLayer
	Transactions *transactionLayer
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		Customers:    newCustomerLayer(db),
		Transactions: newTransactionLayer(db),
	}
}
