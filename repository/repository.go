package repository

import "gorm.io/gorm"

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
