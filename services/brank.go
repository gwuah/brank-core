package services

import (
	"brank/core"
	"brank/core/mq"
	"brank/core/utils"
	"brank/repository"
	"net/http"
)

type brankLayer struct {
	repo   repository.Repo
	config *core.Config
	mq     mq.MQ
}

func newBrankLayer(r repository.Repo, c *core.Config, mq mq.MQ) *brankLayer {
	return &brankLayer{
		repo:   r,
		config: c,
		mq:     mq,
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
