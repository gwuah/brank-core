package internal

import (
	"brank/internal/models"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

func RunSeeds(db *gorm.DB) {

	db.Create(&models.Bank{
		Name:            "Fidelity",
		Url:             "http://google.com",
		HasRestEndpoint: Bool(true),
	})

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

}
