package models

import (
	"brank/core/auth"
	"brank/core/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

var QUE_JOBS_TABLE = `
CREATE TABLE IF NOT EXISTS que_jobs (
	priority    smallint    NOT NULL DEFAULT 100,
	run_at      timestamptz NOT NULL DEFAULT now(),
	job_id      bigserial   NOT NULL,
	job_class   text        NOT NULL,
	args        json        NOT NULL DEFAULT '[]'::json,
	error_count integer     NOT NULL DEFAULT 0,
	last_error  text,
	queue       text        NOT NULL DEFAULT '',
	
	CONSTRAINT que_jobs_pkey PRIMARY KEY (queue, priority, run_at, job_id)
);

COMMENT ON TABLE que_jobs IS '3';
`

func SeedBanks(db *gorm.DB) {
	banks := []Bank{
		{
			Name:            "Standard Chartered",
			Url:             "https://retail.sc.com/afr/ibank/gh/foa/login.htm",
			HasRestEndpoint: utils.Bool(false),
			Code:            "scb",
		},
		{
			Name:            "Fidelity Bank",
			Url:             "https://retailibank.fidelitybank.com.gh/auth/login",
			HasRestEndpoint: utils.Bool(true),
			Code:            "fb",
		},
		{
			Name:            "First National Bank",
			Url:             "https://www.firstnationalbank.com.gh/",
			HasRestEndpoint: utils.Bool(false),
			Code:            "fnb",
		},
	}

	for i := 0; i < len(banks); i++ {
		bank := banks[i]
		if err := db.Model(Bank{}).Where("name=?", bank.Name).First(&bank).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&bank)
			} else {
				log.Println("err is nil", err)
			}
		}
	}
}

func SeedClient(db *gorm.DB) {
	clients := []Client{
		{
			FirstName:   "Mister",
			LastName:    "Brank",
			Email:       "brank@gmail.com",
			Password:    "43d4343i4j3434i44",
			CompanyName: "Brank",
		},
	}

	for i := 0; i < len(clients); i++ {
		client := clients[i]
		if err := db.Model(Client{}).Where("email=?", client.Email).First(&client).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&client)
			} else {
				log.Println("err is nil", err)
			}

		}
	}
}

func SeedApp(db *gorm.DB) {
	apps := []App{
		{
			Name:        "Float",
			ClientID:    1,
			PublicKey:   "934hgreg83r3rv38r3",
			AccessToken: "3B44B34934U30493",
			Logo:        "https://google.com",
			CallbackUrl: "https://google.com",
		},
	}

	for i := 0; i < len(apps); i++ {
		db.Create(&apps[i])
	}
}

func RunSeeds(db *gorm.DB) {

	SeedBanks(db)
	SeedClient(db)
	SeedApp(db)

	db.Create(&Link{
		Code:     auth.GenerateExchangeCode(),
		BankID:   1,
		AppID:    1,
		Username: "banku",
		Password: "stew",
	})

	if err := db.Exec(QUE_JOBS_TABLE).Error; err != nil {
		log.Println("EXEC FAILED")
	}

}
