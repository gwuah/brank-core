package models

type App struct {
	Model
	PublicKey   string `json:"public_key"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	CallbackUrl string `json:"callback_url"`
	AccessToken string `json:"access_token"`
	ClientID    int    `json:"client"`
}
