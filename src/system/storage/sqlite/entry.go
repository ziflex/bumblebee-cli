package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ziflex/bumblebee-ui/src/core"
	"strings"
)

type EntryRepository struct {
	name string
	db   *sql.DB
}

func NewEntryRepository(name string, db *sql.DB) *EntryRepository {
	return &EntryRepository{name, db}
}

func (r *EntryRepository) GetAll() (map[string]*core.Entry, error) {
	rows, err := r.db.Query(fmt.Sprintf("SELECT id, name FROM %s", r.name))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make(map[string]*core.Entry)

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)

		if err != nil {
			break
		}

		results[name] = &core.Entry{
			Id:   id,
			Name: name,
		}
	}

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *EntryRepository) FindOneByName(name string) (*core.Entry, error) {
	if name == "" {
		return nil, errors.New("name must be non-empty string")
	}

	stmt, err := r.db.Prepare(fmt.Sprintf("SELECT id FROM %s where name='?' LIMIT 1", r.name))

	if err != nil {
		return nil, errors.New(err)
	}

	defer stmt.Close()

	var id int

	err = stmt.QueryRow(name).Scan(&id)

	if err != nil {
		return nil, errors.New(err)
	}

	return &core.Entry{Id: id, Name: name}, nil
}

func (r *EntryRepository) Create(entry *core.Entry) error {
	if entry.Id > 0 {
		return errors.Errorf("entry already exists: %s", entry.Name)
	}

	if entry.Name == "" {
		return errors.New("entry must have non-empty name")
	}

	stmt, err := r.db.Prepare(fmt.Sprintf("INSERT INTO %s (name) VALUES ('?')", r.name))

	if err != nil {
		return errors.New(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(entry.Name)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *EntryRepository) CreateMany(entries []*core.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	tx, err := r.db.Begin()

	if err != nil {
		return errors.New(err)
	}

	valueStrings := make([]string, 0, len(entries))
	valueArgs := make([]interface{}, 0, len(entries)*3)

	for _, entry := range entries {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, entry.Name)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (name) VALUES %s",
		r.name,
		strings.Join(valueStrings, ","),
	)

	_, err = tx.Exec(query, valueArgs...)

	if err != nil {
		if rollbackFailure := tx.Rollback(); rollbackFailure != nil {
			err = errors.Errorf("%s \n %s", err.Error(), rollbackFailure.Error())
		}

		return errors.New(err)
	}

	err = tx.Commit()

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *EntryRepository) Delete(entry *core.Entry) error {
	if entry.Id == 0 {
		return errors.Errorf("entry does not exist: %s", entry.Name)
	}

	stmt, err := r.db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id=?", r.name))

	if err != nil {
		return errors.New(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(entry.Id)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r *EntryRepository) DeleteManyByName(names []string) error {
	if len(names) == 0 {
		return nil
	}

	tx, err := r.db.Begin()

	if err != nil {
		return errors.New(err)
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE name IN ('%s')",
		r.name,
		strings.Join(names, "', '"),
	)

	_, err = tx.Exec(query)

	if err != nil {
		if rollbackFailure := tx.Rollback(); rollbackFailure != nil {
			err = errors.Errorf("%s \n %s", err.Error(), rollbackFailure.Error())
		}

		return errors.New(err)
	}

	err = tx.Commit()

	if err != nil {
		return errors.New(err)
	}

	return nil
}
