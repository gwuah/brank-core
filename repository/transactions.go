package repository

import (
	"brank/core/models"
	"math"

	"gorm.io/gorm"
)

type transactionLayer struct {
	db *gorm.DB
}

func newTransactionLayer(db *gorm.DB) *transactionLayer {
	return &transactionLayer{
		db: db,
	}
}

func (t *transactionLayer) BulkInsertWithCount(records *[]models.Transaction, count int) error {
	return t.db.CreateInBatches(records, count).Error
}

func (t *transactionLayer) FindById(id int) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := t.db.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionLayer) Find(query map[string]interface{}, page int) (*TransactionPaginator, error) {
	transactions := []models.Transaction{}

	// seed it with query before passing it down to the paginator
	paginator := Paging(&Param{
		DB:      t.db.Where(query),
		Page:    page,
		Limit:   100,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &transactions)

	return paginator, nil
}

// ignore this portion of the transactions repository
type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

type TransactionPaginator struct {
	TotalRecord int64                `json:"total_record"`
	TotalPage   int                  `json:"total_page"`
	Records     []models.Transaction `json:"records"`
	Offset      int                  `json:"offset"`
	Limit       int                  `json:"limit"`
	Page        int                  `json:"page"`
	PrevPage    int                  `json:"prev_page"`
	NextPage    int                  `json:"next_page"`
}

func Paging(p *Param, result *[]models.Transaction) *TransactionPaginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	var paginator TransactionPaginator
	var count int64
	var offset int

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	db.Model(result).Count(&count)

	paginator.TotalRecord = count
	paginator.Records = *result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}

	if count == 0 {
		paginator.PrevPage = 0
		paginator.NextPage = 0
	}

	return &paginator
}
