package workers

import (
	"brank/core/models"
	"brank/core/queue"
	"brank/core/utils"
	"brank/integrations/fidelity"
	"brank/repository"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bgentry/que-go"
)

const (
	FidelityTransactionsProcessingJob queue.JobIdentifier = "fidelity_transactions_processing_job"
)

type FidelityTransactions struct {
	r repository.Repo
}

func NewFidelityTransactionsWorker(r repository.Repo) *FidelityTransactions {
	return &FidelityTransactions{
		r: r,
	}
}

func CreateFidelityTransactionsJob(appLinkID int) *FidelityTransactionsJobPayload {
	return &FidelityTransactionsJobPayload{
		AppLinkID: appLinkID,
	}
}

type FidelityTransactionsJobPayload struct {
	AppLinkID int `json:"app_link_id"`
}

func (ft *FidelityTransactions) Identifier() queue.JobIdentifier {
	return FidelityTransactionsProcessingJob
}

func (ft *FidelityTransactions) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args FidelityTransactionsJobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("fidelity_transactions_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		appLink, err := ft.r.AppLink.FindById(args.AppLinkID)
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to find app-link. err:%v", err)
		}

		link, err := ft.r.Link.FindById(appLink.LinkID)
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to find link. err:%v", err)
		}

		meta, err := link.GetMeta()
		if err != nil {
			return fmt.Errorf("fidelity_transactions_worker: get link meta failed. err:%v", err)
		}

		for _, tree := range meta.Fidelity.Trees {
			var (
				transactions []models.Transaction
			)

			for _, transaction := range tree.Transactions {
				var (
					direction models.Direction
					amount    float64
				)
				dateChunks := strings.Split(transaction.Date, "-")
				dateTime := time.Date(utils.ConvertToInt(dateChunks[2]), time.Month(utils.ConvertToInt(dateChunks[1])), utils.ConvertToInt(dateChunks[0]), 0, 0, 0, 0, time.UTC)

				if transaction.Credit != "-" {
					direction = models.Credit
					amount = fidelity.ConvertToFloat(transaction.Credit)
				} else {
					direction = models.Debit
					amount = fidelity.ConvertToFloat(transaction.Debit)
				}

				transactions = append(transactions, models.Transaction{
					Date:        dateTime,
					Description: strings.Join(transaction.Description[:], ","),
					Direction:   direction,
					AccountID:   8, // TODO: change this to tree.AccountID
					Amount:      int64(amount),
					Status:      models.Success,
				})
			}

			if len(transactions) > 0 {
				if err := ft.r.Transactions.BulkInsertWithCount(&transactions, 300); err != nil {
					return fmt.Errorf("fidelity_transactions_worker: failed to bullk insert transactions. err:%v", err)
				}
			}
		}

		return nil
	}
}
