package models

type Bank struct {
	Model
	Code            string `json:"code"`
	Name            string `json:"name"`
	Url             string `json:"url"`
	HasRestEndpoint *bool  `json:"has_rest_endpoint"`
}
