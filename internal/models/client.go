package models

type Client struct {
	Model
	AccessToken string `json:"access_token"`
	Verified    *bool  `json:"verified" gorm:"default=false"`
}
