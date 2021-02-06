package repository

import (
	"brank/core/models"

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

func (ca *clientApplicationLayer) Create(app *models.App) error {
	if err := ca.db.Create(app).Error; err != nil {
		return err
	}
	return nil
}

func (ca *clientApplicationLayer) Update(app *models.App) error {
	if err := ca.db.Save(&app).Error; err != nil {
		return err
	}
	return nil
}

func (ca *clientApplicationLayer) FindByPublicKey(key string) (*models.App, error) {
	app := models.App{PublicKey: key}
	if err := ca.db.Where("public_key = ?", key).First(&app).Error; err != nil {
		return &app, err
	}
	return &app, nil
}
