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
	WebhookJob queue.JobIdentifier = "webhook_job"
)

type Webook struct {
}

type WebHookJobPayload struct {
}

func CreateWebhookJob(linkID int) *WebHookJobPayload {
	return &WebHookJobPayload{}
}

func NewWebhookWorker(i *integrations.Integrations, r repository.Repo) *Webook {
	return &Webook{}
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
