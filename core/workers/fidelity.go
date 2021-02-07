package workers

import (
	"brank/core/models"
	"brank/core/queue"
	"brank/integrations"
	"brank/integrations/fidelity"
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

		// TODO: find by phone number before creating
		// create customer
		c := models.Customer{
			Name:        meta.Fidelity.Otp.User.Name,
			PhoneNumber: meta.Fidelity.Otp.User.MobileNumber,
			BankID:      link.BankID,
		}
		if err := f.r.Customer.Create(&c); err != nil {
			return fmt.Errorf("fidelity_worker: failed to create customer. err:%v", err)
		}

		// TODO: find before create, gyimim
		// create accounts
		var accounts []models.Account
		for _, account := range meta.Fidelity.Otp.User.Accounts {
			accounts = append(accounts, models.Account{
				LinkID:        link.ID,
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
				if err := f.r.Account.UpdateWhere(&models.Account{
					Balance: object.Balance,
				}, "external_id=?", object.Id); err != nil {
					return fmt.Errorf("fidelity_worker: failed to update balance. err:%v", err)
				}
			}
		} else {
			return errors.New("if we're here, then i'm pretty sure the bearer token has expired")
		}

		// pull transactions for each account and store it
		for _, account := range accounts {
			status, response, err := f.i.Fidelity.DownloadStatement(account.ExternalID, fidelity.Get1YearFromToday(), fidelity.GetTodaysDate())
			if err != nil {
				return fmt.Errorf("fidelity_worker: failed to download statement. err:%v", err)
			}

			if status {
				tree, err := f.i.Fidelity.ProcessPDF(response)
				if err != nil {
					return fmt.Errorf("fidelity_worker: failed to get process statement. err:%v", err)
				}

				tree.PopulateSummary()
				meta.Fidelity.Trees = append(meta.Fidelity.Trees, *tree)
			} else {
				return errors.New("if we're here, then i'm pretty sure the bearer token has expired")
			}

		}

		if err := link.CommitMeta(meta); err != nil {
			return fmt.Errorf("fidelity_worker: failed to commit link meta. err:%v", err)
		}

		if err := f.r.Link.Update(link); err != nil {
			return fmt.Errorf("fidelity_worker: failed to update link. err:%v", err)
		}

		return nil
	}
}
