package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Link struct {
	Model
	Raw      postgres.Jsonb `json:"raw"`
	Code     string         `json:"code"`
	BankID   int            `json:"bank_id"`
	Bank     Bank           `json:"bank,omitempty"`
	AppID    int            `json:"app_id"`
	App      App            `json:"app,omitempty"`
	Hash     string         `json:"hash"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}
