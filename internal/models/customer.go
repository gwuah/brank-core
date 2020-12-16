package models

type Customer struct {
	Model
	Hash     string `json:"hash"`
	Username string `json:"username"`
	Password string `json:"password"`
	Deleted  *bool  `json:"deleted"`
	Banks    []Bank `json:"banks" gorm:"many2many:accounts;"`
}
