package internal

import (
	"brank/internal/models"
	"errors"
	"log"

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

func SeedClient(db *gorm.DB) {
	clients := []models.Client{
		{
			FirstName:   "Mister",
			LastName:    "Brank",
			Email:       "brank@gmail.com",
			Password:    "43d4343i4j3434i44",
			CompanyName: "Brank",
		},
	}

	for i := 0; i < len(clients); i++ {
		client := clients[i]
		if err := db.Model(models.Client{}).Where("email=?", client.Email).First(&client).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&client)
			} else {
				log.Println("err is nil", err)
			}

		}
	}
}

func SeedApp(db *gorm.DB) {
	apps := []models.App{
		{
			Name:        "Float",
			ClientID:    1,
			PublicKey:   "934hgreg83r3rv38r3",
			AccessToken: "3B44B34934U30493",
			Logo:        "https://google.com",
			CallbackUrl: "https://google.com",
		},
	}

	for i := 0; i < len(apps); i++ {
		db.Create(&apps[i])
	}
}

func RunSeeds(db *gorm.DB) {

	SeedBanks(db)
	SeedClient(db)
	SeedApp(db)

	db.Create(&models.Link{
		Code:     GenerateExchangeCode(),
		BankID:   1,
		AppID:    1,
		Username: "banku",
		Password: "stew",
	})

	// for i := 0; i < 6000; i++ {
	// 	db.Create(&models.Transaction{
	// 		Direction:   "credit",
	// 		Amount:      5000,
	// 		Description: "Incoming Transfer Clearance",
	// 		Date:        time.Now(),
	// 		Status:      "success",
	// 		InquiryID:   1,
	// 		AccountID:   1,
	// 	})
	// }

}
