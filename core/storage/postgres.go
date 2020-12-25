package storage

import (
	"brank/core"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JoinTableConfig struct {
	Model     interface{}
	JoinTable interface{}
	Field     string
}

func SetupJoinTables(db *gorm.DB, configs []JoinTableConfig) error {

	sjt := func(config JoinTableConfig) error {
		err := db.SetupJoinTable(config.Model, config.Field, config.JoinTable)
		return err
	}

	for _, config := range configs {
		if err := sjt(config); err != nil {
			return nil
		}
	}

	return nil
}

func RunMigrations(db *gorm.DB, models ...interface{}) error {
	err := db.AutoMigrate(models...)
	return err
}

func GeneratePostgresURI(config *core.Config) string {
	var (
		dbUrl    = config.DATABASE_URL
		host     = config.PG_HOST
		port     = config.PG_PORT
		dbname   = config.PG_NAME
		user     = config.PG_USER
		password = config.PG_PASS
		sslmode  = config.PG_SSLMODE
	)
	if config.ENVIRONMENT == core.Development {
		dbUrl = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	}
	return dbUrl
}

func NewPostgres(config *core.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	db, err = gorm.Open(postgres.Open(GeneratePostgresURI(config)), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
