package repository

import (
	"brank/internal/models"

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

func (cl *linkLayer) FindByCode(code string) (*models.Link, error) {
	link := models.Link{Code: code}
	if err := cl.db.Where("code = ?", code).First(&link).Error; err != nil {
		return &link, err
	}
	return &link, nil
}

func (r *linkLayer) Create(link *models.Link) error {
	if err := r.db.Create(link).Error; err != nil {
		return err
	}
	return nil
}

func (r *linkLayer) Update(link *models.Link) error {
	if err := r.db.Save(&link).Error; err != nil {
		return err
	}
	return nil
}
