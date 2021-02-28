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
	q *queue.Que
}

type FidelityJobPayload struct {
	AppLinkID int `json:"app_link_id"`
}

func CreateFidelityJob(appLinkID int) *FidelityJobPayload {
	return &FidelityJobPayload{
		AppLinkID: appLinkID,
	}
}

func NewFidelityWorker(i *integrations.Integrations, r repository.Repo, q *queue.Que) *Fidelity {
	return &Fidelity{
		i: i,
		r: r,
		q: q,
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

		appLink, err := f.r.AppLink.FindById(args.AppLinkID)
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to find app-link. err:%v", err)
		}

		link, err := f.r.Link.FindById(appLink.LinkID)
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to find link. err:%v", err)
		}

		meta, err := link.GetMeta()
		if err != nil {
			return fmt.Errorf("fidelity_worker: get link meta failed. err:%v", err)
		}

		// create customer
		customer, err := f.r.Customer.FindByPhone(meta.Fidelity.Otp.User.MobileNumber)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			customer.Name = meta.Fidelity.Otp.User.Name
			customer.PhoneNumber = meta.Fidelity.Otp.User.MobileNumber
			customer.BankID = link.BankID
			if err := f.r.Customer.Create(customer); err != nil {
				return fmt.Errorf("fidelity_worker: failed to create customer. err:%v", err)
			}
		}

		var externalIdToAccountNumber map[int64]string

		// create accounts
		var accounts []models.Account
		for _, acc := range meta.Fidelity.Otp.User.Accounts {
			var tmpAccount models.Account
			externalIdToAccountNumber[acc.Id] = acc.AccountNumber
			if err := f.r.Account.FindWhere(&tmpAccount, "account_number=? AND external_id=?", acc.AccountNumber, acc.Id); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					accounts = append(accounts, models.Account{
						LinkID:        link.ID,
						AccountNumber: acc.AccountNumber,
						Currency:      acc.Currency,
						ExternalID:    acc.Id,
						Name:          acc.Name,
						CustomerID:    customer.ID,
					})
				} else {
					return fmt.Errorf("fidelity_worker: failed to check for account. err:%v", err)
				}
			}
		}
		if len(accounts) > 0 {
			if err := f.r.Account.BulkInsert(&accounts); err != nil {
				return fmt.Errorf("fidelity_worker: failed to create accounts. err:%v", err)
			}
		}

		// get balance
		f.i.Fidelity.SetBearerToken(meta.Fidelity.Otp.Token)
		response, err := f.i.Fidelity.GetBalance()
		if err != nil {
			return fmt.Errorf("fidelity_worker: failed to get balances. err:%v , response: %v", err, response)
		}

		for _, balanceObject := range response.Balances {
			if err := f.r.Account.UpdateWhere(&models.Account{
				Balance: balanceObject.Balance,
			}, "account_number=?", externalIdToAccountNumber[balanceObject.Id]); err != nil {
				return fmt.Errorf("fidelity_worker: failed to update balance. err:%v", err)
			}
		}

		// pull transactions for each account and store it
		f.i.Fidelity.SetBearerToken(meta.Fidelity.Otp.Token)
		// a proper implementation is to load all accounts a
		// you can put a field that shows the last time a sync was done.
		// if it's less than a day, we ignore
		// we only need account.ExternalId, no need to load everything

		// the thing with the current approach is, if the creating account works
		// and the the code below doesn't work. when the task re-tries the job,
		// the accounts wont be seeded
		var accs = &accounts
		if len(accounts) == 0 {
			accs, err = f.r.Account.Find("customer_id=?", customer.ID)
			if err != nil {
				return fmt.Errorf("fidelity_worker: failed to find link accounts. err:%v", err)
			}
		}
		for _, account := range *accs {
			response, err := f.i.Fidelity.DownloadStatement(account.ExternalID, fidelity.Get3YearsFromToday(), fidelity.GetTodaysDate())
			if err != nil {
				return fmt.Errorf("fidelity_worker: failed to download statement. err:%v", err)
			}

			tree, err := f.i.Fidelity.ProcessPDF(response)
			if err != nil {
				return fmt.Errorf("fidelity_worker: failed to get process statement. err:%v", err)
			}

			tree.PopulateSummary()
			tree.AccountID = account.ID
			meta.Fidelity.Trees = append(meta.Fidelity.Trees, *tree)

		}

		if err := link.CommitMeta(meta); err != nil {
			return fmt.Errorf("fidelity_worker: failed to commit link meta. err:%v", err)
		}

		if err := f.r.Link.Update(link); err != nil {
			return fmt.Errorf("fidelity_worker: failed to update link. err:%v", err)
		}

		if err := f.q.QueueJob(FidelityTransactionsProcessingJob, CreateFidelityTransactionsJob(appLink.ID)); err != nil {
			return fmt.Errorf("fidelity_worker: failed to queue processing job. err:%v", err)
		}

		return nil
	}
}
