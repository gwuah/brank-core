package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Inquiry struct {
	Model
	Raw       postgres.Jsonb `json:"raw"`
	AccountID int            `json:"account_id"`
	Account   Account        `json:"account,omitempty"`
}
