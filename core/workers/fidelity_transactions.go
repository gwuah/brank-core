package workers

import (
	"brank/core/queue"
	"brank/integrations"
	"brank/repository"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bgentry/que-go"
)

const (
	FidelityTransactionsJob queue.JobIdentifier = "fidelity_transactions_job"
)

type FidelityTransactions struct {
}

type FidelityTransactionsJobPayload struct {
}

func NewFidelityTransactionsWorker(i *integrations.Integrations, r repository.Repo) *FidelityTransactions {
	return &FidelityTransactions{}
}

func (w *FidelityTransactions) Identifier() queue.JobIdentifier {
	return FidelityTransactionsJob
}

func (w *FidelityTransactions) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args FidelityTransactionsJobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("fidelity_transactions_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		return errors.New("intentional, to retain job in queue")
	}
}
