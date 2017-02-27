package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ziflex/bumblebee-cli/src/core"
)

type SettingsRepository struct {
	name string
	db   *sql.DB
}

func NewSettingsRepository(name string, db *sql.DB) *SettingsRepository {
	return &SettingsRepository{name, db}
}

func (r *SettingsRepository) Get() (*core.Settings, error) {
	rows, err := r.db.Query(fmt.Sprintf("SELECT key, value FROM %s", r.name))

	if err != nil {
		return nil, errors.New(err)
	}

	result := core.NewDefaultSettings()

	for rows.Next() {
		var key string
		var value string

		err = rows.Scan(&key, &value)

		if err != nil {
			break
		}

		switch key {
		case "prefix":
			if value != "" {
				result.Prefix = value
			}
		case "directory":
			if value != "" {
				result.Directory = value
			}
		}
	}

	if err != nil {
		return nil, errors.New(err)
	}

	return result, nil
}

func (r *SettingsRepository) Set(settings *core.Settings) error {
	tx, err := r.db.Begin()

	if err != nil {
		return errors.New(err)
	}

	values := map[string]string{
		"prefix":    settings.Prefix,
		"directory": settings.Directory,
	}

	for key, value := range values {
		err = r.createOrUpdate(tx, key, value)

		if err != nil {
			break
		}
	}

	if err != nil {
		if rollbackFailure := tx.Rollback(); rollbackFailure != nil {
			return errors.Errorf("%s \n %s", err.Error(), rollbackFailure.Error())
		}

		return errors.New(err)
	}

	err = tx.Commit()

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *SettingsRepository) createOrUpdate(tx *sql.Tx, key, value string) error {
	update, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET value=? where key=?", r.name))

	if err != nil {
		return err
	}

	res, err := update.Exec(value, key)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected > 0 {
		return nil
	}

	insert, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (key, value) VALUES(?, ?)", r.name))

	if err != nil {
		return err
	}

	_, err = insert.Exec(key, value)

	return err
}
