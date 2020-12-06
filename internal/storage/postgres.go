package storage

import (
	"brank/internal"
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

func NewPostgres(config *internal.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.PG_HOST, config.PG_PORT, config.PG_USER, config.PG_PASS, config.PG_NAME, config.PG_SSLMODE,
	)

	env := internal.GetEnvironment()
	if env == internal.Staging || env == internal.Production {
		db, err = gorm.Open(postgres.Open(config.DATABASE_URL), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
