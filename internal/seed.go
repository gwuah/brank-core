package internal

import (
	"brank/internal/models"
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
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

	db.Create(&models.Client{
		AccessToken: "439y4y3g7ggr38r3grg37r",
		Verified:    Bool(false),
	})

	db.Create(&models.Customer{
		Hash:    "our38y4834937g7344",
		Deleted: Bool(false),
	})

	db.Create(&models.Account{
		Name:             "Fidelity Lifestyle",
		Balance:          67000,
		AvailableBalance: 66000,
		CustomerID:       1,
		BankID:           1,
		Credentials:      postgres.Jsonb{RawMessage: []byte("")},
		Deleted:          Bool(false),
	})

	db.Create(&models.Inquiry{
		Raw:       postgres.Jsonb{RawMessage: []byte("")},
		AccountID: 1,
	})

	db.Create(&models.Transaction{
		Direction:   "credit",
		Amount:      67000,
		Description: "TransferWise Incoming",
		Date:        time.Now(),
		Status:      "success",
		InquiryID:   1,
		AccountID:   1,
	})

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
