package integrations

import "brank/internal"

type BankIntegration struct {
	axios    internal.Axios
	endpoint string
}

type Integrations struct {
	Fidelity *BankIntegration
}

func NewBankIntegrations() *Integrations {
	a := internal.NewAxios()
	return &Integrations{
		Fidelity: NewFidelity(a),
	}
}
