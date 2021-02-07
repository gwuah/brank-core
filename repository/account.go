package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type accountLayer struct {
	db *gorm.DB
}

func newAccountLayer(db *gorm.DB) *accountLayer {
	return &accountLayer{
		db: db,
	}
}

func (a *accountLayer) BulkInsert(records *[]models.Account) error {
	return a.db.Create(records).Error
}

func (a *accountLayer) Update(records *models.Account) error {
	return a.db.Updates(records).Error
}
