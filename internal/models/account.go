package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Account struct {
	Model
	Name             string         `json:"name"`
	Balance          int64          `json:"balance"`
	AvailableBalance int64          `json:"available_balance"`
	Transactions     []Transaction  `json:"transactions"`
	CustomerID       int            `json:"customer_id"`
	BankID           int            `json:"bank_id"`
	Customer         Customer       `json:"customer,omitempty"`
	Bank             Bank           `json:"bank,omitempty"`
	Credentials      postgres.Jsonb `json:"credentials"`
	Inquiries        []Inquiry      `json:"inquiries"`
	Deleted          *bool          `json:"deleted"`
}
