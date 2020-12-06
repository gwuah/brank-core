package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Inquiry struct {
	Model
	Raw            postgres.Jsonb `json:"raw"`
	CustomerBankID int            `json:"customer_bank_id"`
	CustomerBank   CustomerBank   `json:"customer_bank"`
}
