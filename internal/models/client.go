package models

type Client struct {
	Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
	Verified    *bool  `json:"verified" gorm:"default=false"`
}
