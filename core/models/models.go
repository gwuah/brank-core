package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Direction string
type Status string

var (
	Debit   Direction = "debit"
	Credit  Direction = "credit"
	Failed  Status    = "failed"
	Success Status    = "success"
)

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

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

type App struct {
	Model
	PublicKey   string `json:"public_key"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	CallbackUrl string `json:"callback_url"`
	AccessToken string `json:"access_token"`
	ClientID    int    `json:"client"`
}

type Bank struct {
	Model
	Code            string `json:"code"`
	Name            string `json:"name"`
	Url             string `json:"url"`
	HasRestEndpoint *bool  `json:"has_rest_endpoint"`
}

type Client struct {
	Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
	Verified    *bool  `json:"verified" gorm:"default=false"`
}

type Customer struct {
	Model
	Hash string `json:"hash"`
}

type Inquiry struct {
	Model
	Raw    postgres.Jsonb `json:"raw"`
	LinkID int            `json:"link_id"`
	Link   Link           `json:"link,omitempty"`
}

type Link struct {
	Model
	Raw      postgres.Jsonb `json:"raw"`
	Code     string         `json:"code"`
	BankID   int            `json:"bank_id"`
	Bank     *Bank          `json:"bank,omitempty"`
	AppID    int            `json:"app_id"`
	App      *App           `json:"app,omitempty"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}

type Transaction struct {
	Model
	Direction   Direction `json:"direction"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"transaction_date"`
	Status      Status    `json:"status"`
	InquiryID   int       `json:"inquiry_id"`
	Inquiry     *Inquiry  `json:"inquiry,omitempty"`
	AccountID   int       `json:"account_id"`
	Account     *Account  `json:"account,omitempty"`
}

// type QueJob struct {
// 	ID         int64          `json:"id"`
// 	JobId      int            `json:"job_id"`
// 	JobClass   string         `json:"job_class"`
// 	Args       postgres.Jsonb `json:"args"`
// 	RunAt      time.Time      `json:"run_at"`
// 	Priority   int            `json:"priority" gorm:"default:100"`
// 	Queue      string         `json:"queue"`
// 	Type       string         `json:"type"`
// 	ErrorCount int            `json:"error_count"`
// 	LastError  string         `json:"last_error"`
// }