package models

type Client struct {
	Model
	AccessToken string `json:"access_token"`
	CallbackUrl string `json:"callback_url"`
	Verified    *bool  `json:"verified" gorm:"default=false"`
}
