package models

import (
	"brank/integrations/fidelity"
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Direction string
type Status string
type AppLinkState string

var (
	Debit   Direction = "debit"
	Credit  Direction = "credit"
	Failed  Status    = "failed"
	Success Status    = "success"

	Claimed   AppLinkState = "claimed"
	Unclaimed AppLinkState = "unclaimed"
)

type LinkConfiguration struct {
	Otp     []FormConfig `json:"otp"`
	Initial []FormConfig `json:"initial"`
}

type BankMeta struct {
	FormConfiguration LinkConfiguration `json:"link_configuration"`
}

type Model struct {
	ID        int        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

type Account struct {
	Model
	ExternalID       int64   `json:"external_id"`
	Name             string  `json:"name"`
	AccountNumber    string  `json:"account_number"`
	Balance          float64 `json:"balance"`
	Currency         string  `json:"currency"`
	AvailableBalance int64   `json:"available_balance,omitempty"`
	LinkID           int     `json:"link_id"`
	CustomerID       int     `json:"customer_id"`
}

type App struct {
	Model
	PublicKey   string `json:"public_key"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	CallbackUrl string `json:"callback_url"`
	AccessToken string `json:"access_token"`
	ClientID    int    `json:"client_id"`
}

type FormConfig struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Label       string `json:"label,omitempty"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder,omitempty"`
	Value       string `json:"value,omitempty"`
}

type Bank struct {
	Model
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	RequiresOtp *bool          `json:"requires_otp"`
	Meta        postgres.Jsonb `json:"meta"`
}

type Client struct {
	Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	CompanyName string `json:"company_name"`
	Verified    *bool  `json:"verified" gorm:"default=false"`
}

type Customer struct {
	Model
	Name        string `json:"name"`
	PhoneNumber string `json:"phone"`
	Hash        string `json:"hash"`
	BankID      int    `json:"bank_id"`
}

type Fidelity struct {
	Init  fidelity.HTTPResponse      `json:"init"`
	Otp   fidelity.HTTPResponse      `json:"otp"`
	Trees []fidelity.TransactionTree `json:"transaction_trees"`
}

type LinkMeta struct {
	Fidelity Fidelity `json:"fidelity"`
}

type Link struct {
	Model
	BankID   int            `json:"bank_id"`
	Bank     *Bank          `json:"bank,omitempty"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Meta     postgres.Jsonb `json:"meta"`
}

type AppLink struct {
	Model
	Code        string       `json:"code"`
	AccessToken string       `json:"access_token"`
	AppID       int          `json:"app_id"`
	LinkID      int          `json:"link_id"`
	State       AppLinkState `json:"state"`
}

func (b *Bank) GetMeta() (*BankMeta, error) {
	m := BankMeta{}
	if len(b.Meta.RawMessage) > 0 {
		err := json.Unmarshal(b.Meta.RawMessage, &m)
		if err != nil {
			return nil, err
		}
	}
	return &m, nil
}

func (b *Bank) CommitMeta(m *BankMeta) error {
	converted, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b.Meta = postgres.Jsonb{RawMessage: converted}
	return nil
}

func (b *Link) GetMeta() (*LinkMeta, error) {
	m := LinkMeta{}
	if len(b.Meta.RawMessage) > 0 {
		err := json.Unmarshal(b.Meta.RawMessage, &m)
		if err != nil {
			return nil, err
		}
	}
	return &m, nil
}

func (b *Link) CommitMeta(m *LinkMeta) error {
	converted, err := json.Marshal(m)
	if err != nil {
		return err
	}
	b.Meta = postgres.Jsonb{RawMessage: converted}
	return nil
}

type Transaction struct {
	Model
	Direction   Direction `json:"direction"`
	Amount      int64     `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"transaction_date"`
	Status      Status    `json:"status"`
	AccountID   int       `json:"account_id"`
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
