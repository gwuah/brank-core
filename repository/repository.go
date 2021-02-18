package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Total  int `json:"total"`
}

type BulkLoad struct {
	Pagination
	Records []models.Transaction `json:"records"`
}

type Repo struct {
	Transactions *transactionLayer
	Clients      *clientLayer
	Application  *applicationLayer
	Link         *linkLayer
	Bank         *bankLayer
	Customer     *customerLayer
	Account      *accountLayer
	AppLink      *appLinkLayer
}

func New(db *gorm.DB) Repo {
	return Repo{
		Transactions: newTransactionLayer(db),
		Clients:      newClientLayer(db),
		Application:  newApplicationLayer(db),
		Link:         newLinkLayer(db),
		Bank:         newBankLayer(db),
		Customer:     newCustomerLayer(db),
		Account:      newAccountLayer(db),
		AppLink:      newAppLinkLayer(db),
	}
}
