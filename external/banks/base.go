package external

import "brank/core/utils"

type BankIntegration struct {
	axios    utils.Axios
	endpoint string
}

type Integrations struct {
	Fidelity *BankIntegration
}

func NewBankIntegrations() *Integrations {
	a := utils.NewAxios()
	return &Integrations{
		Fidelity: NewFidelity(a),
	}
}
