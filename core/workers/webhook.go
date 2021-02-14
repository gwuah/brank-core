package workers

import (
	"brank/core/queue"
	"brank/repository"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bgentry/que-go"
)

const (
	WebhookJob queue.JobIdentifier = "webhook_job"
)

type Webook struct {
	r repository.Repo
	q *queue.Que
}

type WebHookJobPayload struct {
	AppLinkID int `json:"app_link_id"`
}

func CreateWebhookJob(appLinkID int) *WebHookJobPayload {
	return &WebHookJobPayload{
		AppLinkID: appLinkID,
	}
}

func NewWebhookWorker(r repository.Repo, q *queue.Que) *Webook {
	return &Webook{
		r: r,
		q: q,
	}
}

func (w *Webook) Identifier() queue.JobIdentifier {
	return WebhookJob
}

func (w *Webook) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args WebHookJobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("webhook_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		return errors.New("intentional, to retain job in queue")
	}
}
