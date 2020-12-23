package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Link struct {
	Model
	Raw      postgres.Jsonb `json:"raw"`
	BankID   int            `json:"bank_id"`
	Bank     Bank           `json:"bank,omitempty"`
	ClientID int            `json:"client_id"`
	Client   Client         `json:"client,omitempty"`
	Hash     string         `json:"hash"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}
