package models

import "github.com/jinzhu/gorm/dialects/postgres"

type Inquiry struct {
	Model
	Raw    postgres.Jsonb `json:"raw"`
	LinkID int            `json:"link_id"`
	Link   Link           `json:"link,omitempty"`
}
