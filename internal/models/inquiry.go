package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Inquiry struct {
	Model
	Raw        postgres.Jsonb `json:"raw"`
	CustomerId int            `json:"customer_id"`
	Customer   Customer       `json:"customer"`
}
