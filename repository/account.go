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
		db: db.Debug(),
	}
}

func (a *accountLayer) BulkInsert(records *[]models.Account) error {
	return a.db.Create(records).Error
}

func (a *accountLayer) UpdateWhere(record *models.Account, query string, params ...interface{}) error {
	return a.db.Where(query, params...).Updates(record).Error
}

func (a *accountLayer) FindWhere(record *models.Account, query string, params ...interface{}) error {
	return a.db.Where(query, params...).First(record).Error
}

func (a *accountLayer) Find(query string, params ...interface{}) (*[]models.Account, error) {
	var accounts []models.Account
	if err := a.db.Where(query, params...).Find(&accounts).Error; err != nil {
		return &accounts, err
	}
	return &accounts, nil
}
