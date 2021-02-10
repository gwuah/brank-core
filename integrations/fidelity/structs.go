package fidelity

import "brank/core/axios"

type Integration struct {
	axios                      *axios.Axios
	loginEndpoint              string
	verifyOtpEndpoint          string
	balanceEndpoint            string
	recentTransactionsEndpoint string
	statementEndpoint          string
}

type Account struct {
	AccountNumber string `json:"accountNumber"`
	Currency      string `json:"currency"`
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Type          int    `json:"type"`
}

type User struct {
	UserUid        string    `json:"userUid"`
	UserType       int       `json:"userType"`
	Status         int       `json:"status"`
	Name           string    `json:"name"`
	MobileNumber   string    `json:"mobileNumber"`
	ActiveChannels []int     `json:"activeChannels"`
	Accounts       []Account `json:"accounts"`
}

type Balance struct {
	Id      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

type HTTPResponse struct {
	// login endpoint
	Token              string `json:"token"`
	SecondFactorMethod string `json:"secondFactorMethod"`

	// api request error
	Error     string `json:"error"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`

	// verify otp endpoint
	RemeberMeToken string `json:"rememberMeToken"`
	User           User   `json:"user"`

	// balance endpoint
	Balances              []Balance `json:"balances"`
	TrackingId            string    `json:"trackingId"`
	TransactionStatusCode int       `json:"transactionStatusCode"`
}

type Transaction struct {
	Date        string   `json:"date"`
	Description []string `json:"description"`
	ValueDate   string   `json:"value_date"`
	Debit       string   `json:"debit"`
	Credit      string   `json:"credit"`
	Balance     string   `json:"balance"`
}

type TransactionTree struct {
	TotalCredits float64       `json:"debits_total"`
	TotalDebits  float64       `json:"credits_total"`
	DebitsCount  int64         `json:"debits_count"`
	CreditsCount int64         `json:"credits_count"`
	Transactions []Transaction `json:"transactions"`
	Summary      []string      `json:"summary"`
	AccountID    int           `json:"account_id"`
}

func (tt *TransactionTree) PopulateSummary() {
	indexes := []int{0, 2, 4, 6}
	for _, v := range indexes {
		switch v {
		case 0:
			tt.DebitsCount = ConvertToInt(tt.Summary[v+1])
		case 2:
			tt.CreditsCount = ConvertToInt(tt.Summary[v+1])
		case 4:
			tt.TotalDebits = ConvertToFloat(tt.Summary[v+1])
		case 6:
			tt.TotalCredits = ConvertToFloat(tt.Summary[v+1])
		}

	}
}
