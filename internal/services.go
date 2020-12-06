package internal

import (
	"brank/internal/repository"
	"log"
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

func HandleGetTransactions(req TransactionsRequest, repo repository.Repo) BrankResponse {
	customer, err := repo.Customer.FindById(req.CustomerId)
	if err != nil {
		log.Println(err)
		return BrankResponse{
			Error: true,
			Code:  http.StatusNotFound,
			Meta: BrankMeta{
				Data:    map[string]interface{}{},
				Message: "Customer not found",
			},
		}
	}
	return BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: BrankMeta{
			Data:    customer,
			Message: "Transactions successfully retrieved",
		},
	}

}
