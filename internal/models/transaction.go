package models

type Transaction struct {
	Model
	Direction   Direction `json:"direction"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Date        string    `json:"transaction_date"`
	Status      Status    `json:"status"`
	InquiryId   int       `json:"inquiry_id"`
	Inquiry     Inquiry   `json:"inquiry"`
}
