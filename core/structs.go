package core

type FidelityTransactionsProcessorQeueJob struct {
	LinkID int `json:"link_id"`
}

type MessageRequest struct {
	Message string `json:"message"`
}

type CreateClientRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"company_name"`
}

type LoginClientRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateAppRequest struct {
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	CallbackUrl string `json:"callback_url"`
	ClientID    int    `json:"client_id"`
}

type UpdateAppRequest struct {
	ID int `json:"id"`
	CreateAppRequest
}

type LinkAccountRequest struct {
	BankID    int    `json:"bank_id"`
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type VerifyOTPRequest struct {
	SessionID string `json:"session_id"`
	Otp       string `json:"otp"`
}

type VerifyLoginsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	BankId   int    `json:"bank_id"`
	ClientId int    `json:"client_id"`
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
}

type BrankMeta struct {
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    string      `json:"message"`
}
type BrankResponse struct {
	Error bool      `json:"error"`
	Code  int       `json:"code"`
	Meta  BrankMeta `json:"meta"`
}

// Auth Endpoints
type AuthParams struct {
	AccessToken  string `json:"access_token"`
	AppLinkToken string `json:"link_token"`
}

// Product Endpoints
type TransactionsRequest struct {
	AuthParams
	Pagination
}

type ExchangeContractCode struct {
	AccessToken string `json:"access_token"`
	Code        string `json:"code"`
}
