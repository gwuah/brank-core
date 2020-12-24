package repository

import (
	"brank/internal/models"

	"gorm.io/gorm"
)

type clientApplicationLayer struct {
	db *gorm.DB
}

func newClientApplicationLayer(db *gorm.DB) *clientApplicationLayer {
	return &clientApplicationLayer{
		db: db,
	}
}

func (r *clientApplicationLayer) Create(app *models.App) error {
	if err := r.db.Create(app).Error; err != nil {
		return err
	}
	return nil
}

func (r *clientApplicationLayer) Update(app *models.App) error {
	if err := r.db.Save(&app).Error; err != nil {
		return err
	}
	return nil
}
