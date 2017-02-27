package storage

import (
	"database/sql"
	"github.com/go-errors/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ziflex/bumblebee-cli/src/core"
)

const (
	ENTRY_TABLE    = "entries"
	SETTINGS_TABLE = "settings"
)

type (
	EntryRepository interface {
		GetAll() (map[string]*core.Entry, error)
		FindOneByName(name string) (*core.Entry, error)
		Create(entry *core.Entry) error
		CreateMany(entries []*core.Entry) error
		Delete(entry *core.Entry) error
		DeleteManyByName(names []string) error
	}

	SettingsRepository interface {
		Get() (*core.Settings, error)
		Set(settings *core.Settings) error
	}
)

func OpenDb(location string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", location)

	if err != nil {
		return nil, errors.New(err)
	}

	return db, nil
}
