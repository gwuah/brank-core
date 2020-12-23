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
			Message: "message received",
		},
	}

}

// func HandleGetTransactions(req TransactionsRequest, repo repository.Repo) BrankResponse {
// 	_, err := repo.Customers.FindById(req.CustomerId)
// 	if err != nil {
// 		log.Println("failed to load transactions", err)
// 		return BrankResponse{
// 			Error: true,
// 			Code:  http.StatusNotFound,
// 			Meta: BrankMeta{
// 				Data:    map[string]interface{}{},
// 				Message: "customer not found",
// 			},
// 		}
// 	}
// 	res, err := repo.Transactions.Find(map[string]interface{}{
// 		// "direction": "debit",
// 	}, req.Page)

// 	if err != nil {
// 		log.Println("failed to load transactions", err)
// 		return BrankResponse{
// 			Error: true,
// 			Code:  http.StatusInternalServerError,
// 			Meta: BrankMeta{
// 				Data:    map[string]interface{}{},
// 				Message: "failed to load transactions",
// 			},
// 		}
// 	}

// 	return BrankResponse{
// 		Error: false,
// 		Code:  http.StatusOK,
// 		Meta: BrankMeta{
// 			Data: res.Records,
// 			Pagination: &BrankPagination{
// 				Count:        res.TotalRecord,
// 				NextPage:     res.NextPage,
// 				CurrentPage:  res.Page,
// 				PreviousPage: res.PrevPage,
// 			},
// 			Message: "transactionsyyy successfully retrieved",
// 		},
// 	}

// }
