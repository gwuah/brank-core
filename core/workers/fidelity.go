package worker

import (
	"brank/core/queue"
	"brank/integrations"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bgentry/que-go"
)

const (
	FidelityJob queue.JobIdentifier = "fidelity_job"
)

type Fidelity struct {
	i *integrations.Integrations
}

type JobPayload struct {
	LinkID int `json:"link_id"`
}

func CreateFidelityJob(linkID int) *JobPayload {
	return &JobPayload{
		LinkID: linkID,
	}
}

func NewFidelityWorker(i *integrations.Integrations) *Fidelity {
	return &Fidelity{
		i: i,
	}
}

func (s *Fidelity) Identifier() queue.JobIdentifier {
	return FidelityJob
}

func (s *Fidelity) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args JobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("fidelity_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		// return nil
		return errors.New("intentional, to retain job in queue")
	}
}
