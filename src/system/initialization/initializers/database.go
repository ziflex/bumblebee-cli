package initializers

import (
	"database/sql"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ziflex/bumblebee-ui/src/system/logging"
	"github.com/ziflex/bumblebee-ui/src/system/storage"
	"github.com/ziflex/bumblebee-ui/src/system/utils"
)

var (
	MSG_ERR_DATABASE             = "failed to initialize database"
	MSG_ERR_DATABASE_TRANSACTION = "failed to initialize transaction for database initialization"
)

type (
	DatabaseInitializer struct {
		logger *logging.Logger
		db     *sql.DB
	}

	tableCreator func(tx *sql.Tx) error
)

func NewDatabaseInitializer(logger *logging.Logger, db *sql.DB) *DatabaseInitializer {
	return &DatabaseInitializer{logger, db}
}

func (init *DatabaseInitializer) Run() error {
	tables, err := init.getTableCreators()

	if err != nil {
		init.logger.Error(MSG_ERR_DATABASE)
		init.logger.Error(utils.ErrorStack(err))
		return err
	}

	// Tables are already created
	if len(tables) == 0 {
		return nil
	}

	tx, err := init.db.Begin()

	if err != nil {
		init.logger.Error(MSG_ERR_DATABASE_TRANSACTION)
		return err
	}

	for name, table := range tables {
		failure := table(tx)

		if failure != nil {
			err = failure
			init.logger.Errorf("failed to create table '%s'", name)
			break
		}

		init.logger.Infof("successfully created new table '%s", name)
	}

	if err != nil {
		if rollbackFailure := tx.Rollback(); rollbackFailure != nil {
			err = errors.Errorf("%s \n %s", err.Error(), rollbackFailure.Error())
		}

		init.logger.Error(MSG_ERR_DATABASE)
		init.logger.Error(utils.ErrorStack(err))
		return err
	}

	err = tx.Commit()

	if err != nil {
		init.logger.Error(MSG_ERR_DATABASE)
		init.logger.Error(utils.ErrorStack(err))
		return err
	}

	init.logger.Info("successfully initialized database")

	return nil
}

func (init *DatabaseInitializer) getTableCreators() (map[string]tableCreator, error) {
	rows, err := init.db.Query("SELECT name FROM sqlite_master WHERE type='table'")

	if err != nil {
		return nil, err
	}

	tables := make(map[string]tableCreator)
	tables[storage.ENTRY_TABLE] = init.createEntriesTable
	tables[storage.SETTINGS_TABLE] = init.createSettingsTable

	for rows.Next() {
		var name string

		err = rows.Scan(&name)

		if err != nil {
			break
		}

		switch name {
		case storage.ENTRY_TABLE:
			delete(tables, storage.ENTRY_TABLE)
		case storage.SETTINGS_TABLE:
			delete(tables, storage.SETTINGS_TABLE)
		}
	}

	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (init *DatabaseInitializer) createEntriesTable(tx *sql.Tx) error {
	_, err := tx.Exec(
		fmt.Sprintf(
			"CREATE TABLE %s (id INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL UNIQUE);",
			storage.ENTRY_TABLE,
		),
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		fmt.Sprintf(
			"CREATE UNIQUE INDEX name_idx on %s (name);",
			storage.ENTRY_TABLE,
		),
	)

	if err != nil {
		return err
	}

	return nil
}

func (init *DatabaseInitializer) createSettingsTable(tx *sql.Tx) error {
	_, err := tx.Exec(
		fmt.Sprintf(
			"CREATE TABLE %s(key TEXT, value TEXT);",
			storage.SETTINGS_TABLE,
		),
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		fmt.Sprintf(
			"CREATE UNIQUE INDEX key_idx on %s (key);",
			storage.SETTINGS_TABLE,
		),
	)

	if err != nil {
		return err
	}

	return nil
}
