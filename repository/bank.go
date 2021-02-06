package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type bankLayer struct {
	db *gorm.DB
}

func newBankLayer(db *gorm.DB) *bankLayer {
	return &bankLayer{
		db: db,
	}
}

func (r *bankLayer) FindById(id int) (*models.Bank, error) {
	bank := models.Bank{Model: models.Model{ID: id}}
	if err := r.db.Find(&bank).Error; err != nil {
		return &bank, err
	}
	return &bank, nil
}

func (r *bankLayer) All() (*[]models.Bank, error) {
	var banks []models.Bank
	if err := r.db.Debug().Find(&banks).Error; err != nil {
		return &banks, err
	}
	return &banks, nil
}
