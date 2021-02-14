package repository

import (
	"brank/core/models"

	"gorm.io/gorm"
)

type clientLayer struct {
	db *gorm.DB
}

func newClientLayer(db *gorm.DB) *clientLayer {
	return &clientLayer{
		db: db,
	}
}

func (cl *clientLayer) Create(client *models.Client) error {
	if err := cl.db.Create(client).Error; err != nil {
		return err
	}
	return nil

}

func (cl *clientLayer) FindByEmail(email string) (*models.Client, error) {
	client := models.Client{Email: email}
	if err := cl.db.Where("email = ?", email).First(&client).Error; err != nil {
		return &client, err
	}
	return &client, nil
}

func (cl *clientLayer) FindByID(id int) (*models.Client, error) {
	var client models.Client
	if err := cl.db.First(&client, id).Error; err != nil {
		return &client, err
	}
	return &client, nil
}
