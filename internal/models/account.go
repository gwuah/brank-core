package models

type Account struct {
	Model
	Name             string        `json:"name"`
	Balance          int64         `json:"balance"`
	AvailableBalance int64         `json:"available_balance"`
	Transactions     []Transaction `json:"transactions"`
	LinkID           int           `json:"link_id"`
	Link             Link          `json:"link,omitempty"`
	Deleted          *bool         `json:"deleted"`
}
