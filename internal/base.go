package internal

type MessageRequest struct {
	Message string `json:"message"`
}

type VerifyLoginsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	BankId   int    `json:"bank_id"`
	ClientId int    `json:"client_id"`
}
