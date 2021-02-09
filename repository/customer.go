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
	var customer models.Customer
	if err := cl.db.Where("phone_number=?", phone).First(&customer).Error; err != nil {
		return &customer, err
	}
	return &customer, nil
}

func (cl *customerLayer) Create(customer *models.Customer) error {
	return cl.db.Create(customer).Error
}
