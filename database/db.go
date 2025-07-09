package database

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var DB *sqlx.DB

func InitDBAndMigrate(conStr string) {
	var err error
	DB, err = sqlx.Open("postgres", conStr)
	if err != nil {
		logrus.Errorf("failed to connect to postgres: %v", err)
		return
	}
	err = DB.Ping()
	if err != nil {
		logrus.Errorf("failed to ping postgres DB: %v", err)
		return
	} else {
		logrus.Info("connected to postgres")
	}
	RunningMigration(DB)
}
func RunningMigration(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Errorf("error to create driver %v", err)
		return
	}
	mg, err := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)
	if err != nil {
		logrus.Errorf("error to create the instance %v", err)
		return
	}
	err = mg.Up()
	if err != nil {
		logrus.Errorf("error to running migrations %v", err)
		return
	} else {
		logrus.Info("migrations completed successfully")
		return
	}
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		logrus.Errorf("failed to close DB:%v", err)
		return
	}
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
