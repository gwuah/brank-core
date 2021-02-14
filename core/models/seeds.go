package models

import (
	"brank/core/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

var (
	FidelityBank         = "fb"
	StandardChateredBank = "scb"
	FirstNationalBank    = "fnb"
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

func SeedBanks(db *gorm.DB) error {
	banks := []Bank{
		{
			Name:        "Standard Chartered",
			Code:        StandardChateredBank,
			RequiresOtp: utils.Bool(false),
		},
		{
			Name:        "Fidelity Bank",
			Code:        FidelityBank,
			RequiresOtp: utils.Bool(true),
		},
		{
			Name:        "First National Bank",
			Code:        FirstNationalBank,
			RequiresOtp: utils.Bool(false),
		},
	}

	meta, err := banks[1].GetMeta()
	if err != nil {
		log.Println("failed to get meta", err)
		return err
	}

	meta.FormConfiguration = LinkConfiguration{
		Initial: []FormConfig{
			{
				ID:          1,
				Label:       "Phone Number",
				Type:        "text",
				Required:    true,
				Placeholder: "0205428811",
			},
			{
				ID:          2,
				Label:       "Pin",
				Type:        "text",
				Required:    true,
				Placeholder: "******",
			},
			{
				ID:       3,
				Type:     "button",
				Required: false,
				Value:    "Link Account",
			},
		},
		Otp: []FormConfig{
			{
				ID:          1,
				Label:       "Otp",
				Type:        "text",
				Required:    true,
				Placeholder: "****",
			},
			{
				ID:       2,
				Type:     "button",
				Required: false,
				Value:    "Complete Process",
			},
		},
	}

	if err := banks[1].CommitMeta(meta); err != nil {
		log.Println("failed to commit config meta")
		return err

	}

	for i := 0; i < len(banks); i++ {
		bank := banks[i]
		if err := db.Model(Bank{}).Where("code=?", bank.Code).First(&bank).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&bank).Error; err != nil {
					return err
				}
				continue
			}
			return err
		}
	}

	return nil
}

func SeedClient(db *gorm.DB) error {
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
				if err := db.Create(&client).Error; err != nil {
					return err
				}
				continue
			}
			return err
		}
	}

	return nil
}

func SeedApp(db *gorm.DB) error {
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
		app := apps[i]

		if err := db.Model(App{}).Where("public_key=?", app.PublicKey).First(&app).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&app).Error; err != nil {
					return err
				}
				continue
			}

			return err

		}
	}

	return nil
}

func SeedLink(db *gorm.DB) error {
	var links []Link
	if err := db.Model(Link{}).Find(&links).Limit(1).Error; err != nil {
		return err
	}
	if len(links) > 0 {
		return nil
	}
	return db.Create(&Link{
		BankID:   1,
		Username: "banku",
		Password: "stew",
	}).Error
}

func SeedAppLink(db *gorm.DB) error {
	var appLink []AppLink
	if err := db.Model(AppLink{}).Find(&appLink).Limit(1).Error; err != nil {
		return err
	}
	if len(appLink) > 0 {
		return nil
	}
	return db.Create(&AppLink{
		Code:  utils.GenerateExchangeCode(),
		AppID: 1,
	}).Error
}

func RunSeeds(db *gorm.DB) {

	if err := SeedBanks(db); err != nil {
		log.Println("failed to seed banks")
	}

	if err := SeedClient(db); err != nil {
		log.Println("failed to seed clients")
	}

	if err := SeedApp(db); err != nil {
		log.Println("failed to seed apps")
	}

	if err := SeedLink(db); err != nil {
		log.Println("failed to seed links")
	}

	if err := db.Exec(QUE_JOBS_TABLE).Error; err != nil {
		log.Println("failed to exec queue jobs table ")
	}

}
