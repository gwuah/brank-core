package services

import (
	"brank/core"
	"brank/core/mq"
	"brank/core/queue"
	"brank/core/utils"
	worker "brank/core/workers"
	"brank/repository"

	"net/http"
)

type brankLayer struct {
	repo   repository.Repo
	config *core.Config
	mq     mq.MQ
	q      *queue.Que
}

func newBrankLayer(r repository.Repo, c *core.Config, mq mq.MQ, q *queue.Que) *brankLayer {
	return &brankLayer{
		repo:   r,
		config: c,
		mq:     mq,
		q:      q,
	}
}

func (b *brankLayer) PublishMessageIntoKafka(req core.MessageRequest) core.BrankResponse {
	go func() {
		b.mq.Publish(utils.GenerateTopic("validate_login"), []byte(req.Message))
	}()

	return core.BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: core.BrankMeta{
			Data:    nil,
			Message: "message received",
		},
	}

}

func (b *brankLayer) CreateFidelityTransactionProcessorJob(req core.FidelityTransactionsProcessorQeueJob) core.BrankResponse {

	if err := b.q.QueueJob(worker.FidelityTransactionsProcessingJob, worker.CreateFidelityTransactionsJob(req.LinkID)); err != nil {
		return utils.Error(err, nil, http.StatusInternalServerError)
	}

	return utils.Success(nil, nil)

}

func (b *brankLayer) GetFinancialInstitutions() core.BrankResponse {

	banks, err := b.repo.Bank.All()

	if err != nil {
		return utils.Error(nil, nil, http.StatusInternalServerError)
	}

	return utils.Success(&map[string]interface{}{
		"institutions": banks,
	}, utils.String("financial institutions successfully retrieved"))

}
