package models

import "time"

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

type Direction string
type Status string

var (
	Debit   Direction = "debit"
	Credit            = "credit"
	Failed  Status    = "failed"
	Success           = "success"
)
