package core

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

type CreateAppRequest struct {
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	CallbackUrl string `json:"callback_url"`
	ClientID    int    `json:"client_id"`
}

type LinkAccountRequest struct {
	BankID    string `json:"bank_id"`
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type ExchangeContractCode struct {
	Code string `json:"code"`
}

type TransactionsRequest struct {
	CustomerId int `json:"customer_id"`
	Page       int `json:"page"`
}

type VerifyLoginsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	BankId   int    `json:"bank_id"`
	ClientId int    `json:"client_id"`
}

type BrankPagination struct {
	CurrentPage  int   `json:"current_page,omitempty"`
	NextPage     int   `json:"next_page,omitempty"`
	PreviousPage int   `json:"previous_page,omitempty"`
	Count        int64 `json:"count"`
}

type BrankMeta struct {
	Data       interface{}      `json:"data,omitempty"`
	Pagination *BrankPagination `json:"pagination,omitempty"`
	Message    string           `json:"message"`
}
type BrankResponse struct {
	Error bool      `json:"error"`
	Code  int       `json:"code"`
	Meta  BrankMeta `json:"meta"`
}
