package services

import (
	"brank/core"
	"brank/repository"
	"net/http"
)

type transactionsLayer struct {
	repo   repository.Repo
	config *core.Config
}

func newTransactionLayer(r repository.Repo, c *core.Config) *transactionsLayer {
	return &transactionsLayer{
		repo:   r,
		config: c,
	}
}

func (t *transactionsLayer) GetTransactions(req core.TransactionsRequest) core.BrankResponse {
	// _, err := t.repo.C.FindById(req.CustomerId)
	// if err != nil {
	// 	log.Println("failed to load transactions", err)
	// 	return BrankResponse{
	// 		Error: true,
	// 		Code:  http.StatusNotFound,
	// 		Meta: BrankMeta{
	// 			Data:    map[string]interface{}{},
	// 			Message: "customer not found",
	// 		},
	// 	}
	// }
	// res, err := repo.Transactions.Find(map[string]interface{}{
	// 	// "direction": "debit",
	// }, req.Page)

	// if err != nil {
	// 	log.Println("failed to load transactions", err)
	// 	return BrankResponse{
	// 		Error: true,
	// 		Code:  http.StatusInternalServerError,
	// 		Meta: BrankMeta{
	// 			Data:    map[string]interface{}{},
	// 			Message: "failed to load transactions",
	// 		},
	// 	}
	// }

	return core.BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: core.BrankMeta{
			// Data: res.Records,
			// Pagination: &core.BrankPagination{
			// 	Count:        res.TotalRecord,
			// 	NextPage:     res.NextPage,
			// 	CurrentPage:  res.Page,
			// 	PreviousPage: res.PrevPage,
			// },
			Message: "transactionsyyy successfully retrieved",
		},
	}

}
