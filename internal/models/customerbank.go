package models

import "github.com/jinzhu/gorm/dialects/postgres"

type CustomerBank struct {
	Model
	CustomerId  int            `json:"customer_id"`
	BankId      int            `json:"bank_id"`
	Customer    Customer       `json:"customer"`
	Bank        Bank           `json:"bank"`
	Credentials postgres.Jsonb `json:"credentials"`
	Deleted     *bool          `json:"deleted"`
}
