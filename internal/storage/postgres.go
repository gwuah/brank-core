package storage

import (
	"brank/internal"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
