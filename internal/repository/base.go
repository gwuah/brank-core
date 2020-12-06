package repository

import "gorm.io/gorm"

type Repo struct {
	Customer *customerLayer
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{
		Customer: newCustomer(db),
	}
}
