package internal

import (
	"brank/internal/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedBanks(db *gorm.DB) {
	banks := []models.Bank{
		{
			Name:            "Standard Chartered",
			Url:             "https://retail.sc.com/afr/ibank/gh/foa/login.htm",
			HasRestEndpoint: Bool(false),
			Code:            "scb",
		},
		{
			Name:            "Fidelity Bank",
			Url:             "https://retailibank.fidelitybank.com.gh/auth/login",
			HasRestEndpoint: Bool(true),
			Code:            "fb",
		},
		{
			Name:            "First National Bank",
			Url:             "https://www.firstnationalbank.com.gh/",
			HasRestEndpoint: Bool(false),
			Code:            "fnb",
		},
	}

	for i := 0; i < len(banks); i++ {
		bank := banks[i]
		if err := db.Model(models.Bank{}).Where("name=?", bank.Name).First(&bank).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&bank)
			} else {
				log.Println("err is nil", err)
			}
		}
	}
}

func RunSeeds(db *gorm.DB) {

	SeedBanks(db)

	for i := 0; i < 6000; i++ {
		db.Create(&models.Transaction{
			Direction:   "credit",
			Amount:      5000,
			Description: "Incoming Transfer Clearance",
			Date:        time.Now(),
			Status:      "success",
			InquiryID:   1,
			AccountID:   1,
		})
	}

}
