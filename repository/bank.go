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

func (b *bankLayer) FindById(id int) (*models.Bank, error) {
	bank := models.Bank{Model: models.Model{ID: id}}
	if err := b.db.Find(&bank).Error; err != nil {
		return &bank, err
	}
	return &bank, nil
}

func (b *bankLayer) All() (*[]models.Bank, error) {
	var banks []models.Bank
	if err := b.db.Debug().Find(&banks).Error; err != nil {
		return &banks, err
	}
	return &banks, nil
}
