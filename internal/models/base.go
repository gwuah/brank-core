package models

import "time"

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

type Direction string
type Status string

var (
	Debit   Direction = "debit"
	Credit            = "credit"
	Failed  Status    = "failed"
	Success           = "success"
)
