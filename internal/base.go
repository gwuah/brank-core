package internal

type MessageRequest struct {
	Message string `json:"message"`
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
	CurrentPage int `json:"current_page"`
	NextPage    int `json:"next_page"`
	Count       int `json:"count"`
}

type BrankMeta struct {
	Data       interface{}      `json:"data"`
	Pagination *BrankPagination `json:"pagination,omitempty"`
	Message    string           `json:"message"`
}
type BrankResponse struct {
	Error bool      `json:"error"`
	Code  int       `json:"code"`
	Meta  BrankMeta `json:"meta"`
}
