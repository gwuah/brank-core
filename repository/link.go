package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type linkLayer struct {
	db *gorm.DB
}

func newLinkLayer(db *gorm.DB) *linkLayer {
	return &linkLayer{
		db: db,
	}
}

func (l *linkLayer) FindByCode(code string) (*models.Link, error) {
	link := models.Link{Code: code}
	if err := l.db.Where("code = ?", code).First(&link).Error; err != nil {
		return &link, err
	}
	return &link, nil
}

func (l *linkLayer) Create(link *models.Link) error {
	return l.db.Create(link).Error
}

func (l *linkLayer) Update(link *models.Link) error {
	return l.db.Save(&link).Error
}

func (l *linkLayer) FindById(id int) (*models.Link, error) {
	var link models.Link
	if err := l.db.First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}
