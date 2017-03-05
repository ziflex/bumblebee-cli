package fs

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-ini/ini"
	"path/filepath"
	"strings"
)

type File struct {
	path     string
	engine   *ini.File
	sections map[string]*ini.Section
}

func LoadFile(path string) (*File, error) {
	f, err := ini.Load(path)

	if err != nil {
		return nil, errors.New(err)
	}

	allSections := f.Sections()
	sections := make(map[string]*ini.Section)

	for _, section := range allSections {
		if section.HasKey("Exec") {
			sections[section.Name()] = section
		}
	}

	return &File{
		path:     path,
		engine:   f,
		sections: sections,
	}, nil
}

func (f *File) Name() string {
	return strings.Replace(filepath.Base(f.path), filepath.Ext(f.path), "", 1)
}

func (f *File) GetValues() (map[string]string, error) {
	result := make(map[string]string)
	var err error

	for name, section := range f.sections {
		key, failure := f.getKey(section, "Exec")

		if failure != nil {
			err = failure
			break
		}

		result[name] = key.Value()
	}

	if err != nil {
		return nil, errors.New(err)
	}

	return result, nil
}

func (f *File) SetValues(values map[string]string) error {
	var err error

	for name, value := range values {
		key, failure := f.getKey(f.sections[name], "Exec")

		if failure != nil {
			err = failure
			break
		}

		if key == nil {
			continue
		}

		key.SetValue(value)
	}

	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (f *File) Save() error {
	f.normalizeValues()

	if err := f.engine.SaveTo(f.path); err != nil {
		return errors.New(err)
	}

	return nil
}

func (f *File) normalizeValues() {
	for _, section := range f.sections {
		keys := section.Keys()

		for _, key := range keys {
			val := key.Value()

			if key.Comment != "" && strings.Contains(key.Comment, ";") {
				key.SetValue(f.encode(fmt.Sprintf("%s%s", val, key.Comment)))
				key.Comment = ""
			}
		}
	}
}

func (f *File) getKey(section *ini.Section, keyName string) (*ini.Key, error) {
	if section == nil {
		return nil, nil
	}

	key, err := section.GetKey(keyName)

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (f *File) encode(s string) string {
	return fmt.Sprintf("`%s`", s)
}
