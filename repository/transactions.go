package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type transactionLayer struct {
	db *gorm.DB
}

func newTransactionLayer(db *gorm.DB) *transactionLayer {
	return &transactionLayer{
		db: db.Debug(),
	}
}

func (t *transactionLayer) BulkInsertWithCount(records *[]models.Transaction, count int) error {
	return t.db.CreateInBatches(records, count).Error
}

func (t *transactionLayer) FindById(id int) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := t.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionLayer) Find(p Pagination, query string, params ...interface{}) (*BulkLoad, error) {
	var transactions []models.Transaction
	var total int64

	if err := t.db.Where(query).Limit(p.Limit).Offset(p.Offset).Order("id DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	if err := t.db.Model(transactions).Where(query).Count(&total).Error; err != nil {
		return nil, err
	}

	return &BulkLoad{
		Records: transactions,
		Pagination: Pagination{
			Offset: p.Offset,
			Limit:  p.Limit,
			Total:  total,
		},
	}, nil
}
