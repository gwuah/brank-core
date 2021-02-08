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
		db: db.Debug(),
	}
}

func (cl *customerLayer) FindByPhone(phone string) (*models.Customer, error) {
	customer := models.Customer{PhoneNumber: phone}
	if err := cl.db.First(&customer).Error; err != nil {
		return &customer, err
	}
	return &customer, nil
}

func (cl *customerLayer) Create(customer *models.Customer) error {
	return cl.db.Create(customer).Error
}
