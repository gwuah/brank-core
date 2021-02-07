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
	"gorm.io/gorm"
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
		c, err := f.r.Customer.FindByPhone(meta.Fidelity.Otp.User.MobileNumber)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Name = meta.Fidelity.Otp.User.Name
			c.PhoneNumber = meta.Fidelity.Otp.User.MobileNumber
			c.BankID = link.BankID
			if err := f.r.Customer.Create(c); err != nil {
				return fmt.Errorf("fidelity_worker: failed to create customer. err:%v", err)
			}
		}

		// create accounts
		var accounts []models.Account
		for _, acc := range meta.Fidelity.Otp.User.Accounts {
			var tmpAccount models.Account
			err := f.r.Account.FindWhere(&tmpAccount, "account_number=? AND external_id=?", acc.AccountNumber, acc.Id)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				accounts = append(accounts, models.Account{
					LinkID:        link.ID,
					AccountNumber: acc.AccountNumber,
					Currency:      acc.Currency,
					ExternalID:    acc.Id,
					Name:          acc.Name,
					CustomerID:    c.ID,
				})
			}

			if err != nil {
				return fmt.Errorf("fidelity_worker: failed to check for account. err:%v", err)
			}
		}
		if len(accounts) > 0 {
			if err := f.r.Account.BulkInsert(&accounts); err != nil {
				return fmt.Errorf("fidelity_worker: failed to create accounts. err:%v", err)
			}
		}

		// get balance
		f.i.Fidelity.SetBearerToken(meta.Fidelity.Otp.Token)
		status, response, err := f.i.Fidelity.GetBalance()
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to get balances. err:%v", err)
		}

		if status {
			for _, object := range response.Balances {
				err := f.r.Account.UpdateWhere(&models.Account{
					Balance: object.Balance,
				}, "external_id=?", object.Id)

				if err != nil {
					return fmt.Errorf("fidelity_worker: failed to update balance. err:%v", err)
				}
			}
		} else {
			return errors.New("if we're here, then i'm pretty sure the bearer token has expired")
		}

		// pull transactions for each account and store it
		f.i.Fidelity.SetBearerToken(meta.Fidelity.Otp.Token)
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
