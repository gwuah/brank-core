package workers

import (
	"brank/core/models"
	"brank/core/queue"
	"brank/integrations"
	"brank/repository"
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
	r repository.Repo
}

type FidelityJobPayload struct {
	LinkID int `json:"link_id"`
}

func CreateFidelityJob(linkID int) *FidelityJobPayload {
	return &FidelityJobPayload{
		LinkID: linkID,
	}
}

func NewFidelityWorker(i *integrations.Integrations, r repository.Repo) *Fidelity {
	return &Fidelity{
		i: i,
		r: r,
	}
}

func (f *Fidelity) Identifier() queue.JobIdentifier {
	return FidelityJob
}

func (f *Fidelity) Worker() que.WorkFunc {
	return func(j *que.Job) error {
		var args FidelityJobPayload
		if err := json.Unmarshal(j.Args, &args); err != nil {
			return fmt.Errorf("fidelity_worker: unable to unmarshal job arguments: %v %v", string(j.Args), err)
		}

		link, err := f.r.Link.FindById(args.LinkID)
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to find link. err:%v", err)
		}

		meta, err := link.GetMeta()
		if err != nil {
			return fmt.Errorf("fidelity_worker: get link meta failed. err:%v", err)
		}

		// create customer
		c := models.Customer{
			Name:        meta.Fidelity.Otp.User.Name,
			PhoneNumber: meta.Fidelity.Otp.User.MobileNumber,
			BankID:      link.BankID,
		}
		if err := f.r.Customer.Create(&c); err != nil {
			return fmt.Errorf("fidelity_worker: failed to create customer. err:%v", err)
		}

		// create accounts
		var accounts []models.Account
		for _, account := range meta.Fidelity.Otp.User.Accounts {
			accounts = append(accounts, models.Account{
				AccountNumber: account.AccountNumber,
				Currency:      account.Currency,
				ExternalID:    account.Id,
				Name:          account.Name,
				CustomerID:    c.ID,
			})
		}
		if err := f.r.Account.BulkInsert(&accounts); err != nil {
			return fmt.Errorf("fidelity_worker: failed to create accounts. err:%v", err)
		}

		// get balance
		f.i.Fidelity.SetBearerToken(meta.Fidelity.Otp.Token)
		status, response, err := f.i.Fidelity.GetBalance()
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to get balances. err:%v", err)
		}

		if status {
			for _, object := range response.Balances {
				if err := f.r.Account.Update(&models.Account{
					ExternalID: object.Id,
					Balance:    object.Balance,
				}); err != nil {
					return fmt.Errorf("fidelity_worker: failed to update balance. err:%v", err)
				}
			}
		}

		// pull transactions and seed them
		return errors.New("if we're here, then i'm pretty sure the bearer token has expired")
	}
}
