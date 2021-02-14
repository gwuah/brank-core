package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type appLinkLayer struct {
	db *gorm.DB
}

func newAppLinkLayer(db *gorm.DB) *appLinkLayer {
	return &appLinkLayer{
		db: db,
	}
}

func (l *appLinkLayer) FindByCode(code string) (*models.AppLink, error) {
	appLink := models.AppLink{Code: code}
	if err := l.db.Where("code = ? and state =?", code, models.Unclaimed).First(&appLink).Error; err != nil {
		return &appLink, err
	}
	return &appLink, nil
}

func (l *appLinkLayer) Create(appLink *models.AppLink) error {
	return l.db.Create(appLink).Error
}

func (l *appLinkLayer) Update(appLink *models.AppLink) error {
	return l.db.Debug().Updates(appLink).Error
}

func (l *appLinkLayer) FindById(id int) (*models.AppLink, error) {
	var appLink models.AppLink
	if err := l.db.First(&appLink, id).Error; err != nil {
		return nil, err
	}
	return &appLink, nil
}
