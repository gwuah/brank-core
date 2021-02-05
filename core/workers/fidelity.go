package worker

import (
	"brank/core/queue"
	"encoding/json"
	"fmt"

	"github.com/bgentry/que-go"
)

type SlackJobType string

type Fidelity struct {
}

type JobPayload struct {
}

func NewFidelityWorker() *Fidelity {
	return &Fidelity{}
}

func (s *Fidelity) Identifier() queue.JobIdentifier {
	return queue.FidelityJob
}

func (s *Fidelity) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args JobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("fidelity_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		return nil
	}
}
