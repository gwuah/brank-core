package models

type Customer struct {
	Model
	Hash    string `json:"hash"`
	Deleted *bool  `json:"deleted"`
}
