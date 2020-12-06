package models

import "github.com/jinzhu/gorm/dialects/postgres"

type CustomerBank struct {
	Model
	CustomerID  int            `json:"customer_id"`
	BankID      int            `json:"bank_id"`
	Customer    Customer       `json:"customer"`
	Bank        Bank           `json:"bank"`
	Credentials postgres.Jsonb `json:"credentials"`
	Inquiries   []Inquiry      `json:"inquiries"`
	Deleted     *bool          `json:"deleted"`
}
