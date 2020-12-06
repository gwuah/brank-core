package internal

import (
	"net/http"
)

func HandleMessagePost(req MessageRequest, e EventStore) BrankResponse {
	go func() {
		e.Publish(GenerateTopic("validate_login"), []byte(req.Message))
	}()

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    map[string]interface{}{},
			Message: "Message Received",
		},
	}

}

func HandleGetTransactions(req TransactionsRequest) BrankResponse {

	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    map[string]interface{}{},
			Message: "Transactions successfully retrieved",
			Pagination: &BrankPagination{
				CurrentPage: 1,
				NextPage:    2,
				Count:       100,
			},
		},
	}

}
