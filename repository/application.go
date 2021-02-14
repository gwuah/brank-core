package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type applicationLayer struct {
	db *gorm.DB
}

func newApplicationLayer(db *gorm.DB) *applicationLayer {
	return &applicationLayer{
		db: db,
	}
}

func (a *applicationLayer) Create(app *models.App) error {
	if err := a.db.Create(app).Error; err != nil {
		return err
	}
	return nil
}

func (a *applicationLayer) Update(app *models.App) error {
	if err := a.db.Save(&app).Error; err != nil {
		return err
	}
	return nil
}

func (a *applicationLayer) FindByPublicKey(key string) (*models.App, error) {
	app := models.App{PublicKey: key}
	if err := a.db.Where("public_key = ?", key).First(&app).Error; err != nil {
		return &app, err
	}
	return &app, nil
}
