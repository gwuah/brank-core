package models

type Bank struct {
	Model
	Name string `json:"name"`
	Url  string `json:"url"`
}
