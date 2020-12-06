package models

import "time"

type Transaction struct {
	Model
	Direction   Direction `json:"direction"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"transaction_date"`
	Status      Status    `json:"status"`
	InquiryID   int       `json:"inquiry_id"`
	Inquiry     *Inquiry  `json:"inquiry,omitempty"`
	AccountID   int       `json:"account_id"`
	Account     *Account  `json:"account,omitempty"`
}
