package integrations

import (
	"brank/core/axios"
	"brank/integrations/fidelity"
	"net/http"
	"time"
)

type Integrations struct {
	Fidelity *fidelity.Integration
}

func NewBankIntegrations() *Integrations {
	httpClient := http.Client{Timeout: 30 * time.Second}
	return &Integrations{
		Fidelity: fidelity.New(axios.New(httpClient)),
	}
}
