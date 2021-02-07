package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type customerLayer struct {
	db *gorm.DB
}

func newCustomerLayer(db *gorm.DB) *customerLayer {
	return &customerLayer{
		db: db,
	}
}

func (cl *customerLayer) Create(customer *models.Customer) error {
	return cl.db.Create(customer).Error
}
