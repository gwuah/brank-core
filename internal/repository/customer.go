package repository

import (
	"brank/internal/models"

	"gorm.io/gorm"
)

type customerLayer struct {
	db *gorm.DB
}

func newCustomer(db *gorm.DB) *customerLayer {
	return &customerLayer{
		db: db,
	}
}

func (c *customerLayer) FindById(id int) (*models.Customer, error) {
	customer := &models.Customer{}
	if err := c.db.First(customer, id).Error; err != nil {
		return nil, err
	}
	return customer, nil
}
