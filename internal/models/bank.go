package models

type Bank struct {
	Model
	Name            string     `json:"name"`
	Url             string     `json:"url"`
	HasRestEndpoint *bool      `json:"has_rest_endpoint"`
	Customers       []Customer `json:"customers" gorm:"many2many:accounts;"`
}
