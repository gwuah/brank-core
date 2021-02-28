package services

import (
	"brank/core"
	"brank/core/utils"
	"brank/repository"
	"fmt"
	"net/http"
	"strings"
)

type GetTransactionsParams struct {
	AppLinkID int `json:"app_link_id"`
	Offset    int `json:"offset"`
	Limit     int `json:"limit"`
}

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

func (t *transactionsLayer) GetTransactions(req GetTransactionsParams) core.BrankResponse {
	appLink, err := t.repo.AppLink.FindById(req.AppLinkID)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	accounts, err := t.repo.Account.Find("link_id=?", appLink.LinkID)
	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	var queryArray []string
	for _, account := range *accounts {
		queryArray = append(queryArray, fmt.Sprintf("account_id=%d", account.ID))
	}

	paginationConfig := repository.ValidatePaginationConfig(repository.Pagination{
		Offset: req.Offset,
		Limit:  req.Limit,
	})
	res, err := t.repo.Transactions.Find(paginationConfig, strings.Join(queryArray, " OR "))

	if err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"transactions": res.Records,
		"pagination":   res.Pagination,
	}, nil)

}
