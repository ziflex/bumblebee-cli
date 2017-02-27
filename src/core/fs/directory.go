package fs

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/ziflex/bumblebee-ui/src/system/logging"
	"github.com/ziflex/bumblebee-ui/src/system/utils"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"strings"
)

type Directory struct {
	logger *logging.Logger
	path   string
}

func NewDirectory(logger *logging.Logger, path string) *Directory {
	return &Directory{logger, path}
}

func (d *Directory) Path() string {
	return d.path
}

func (d *Directory) List() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(d.path, "*.desktop"))

	if err != nil {
		return nil, errors.New(err)
	}

	results := make([]string, len(files))

	for i, file := range files {
		name := strings.Replace(file, d.Path()+"/", "", -1)
		name = strings.Replace(name, ".desktop", "", -1)

		results[i] = name
	}

	return results, nil
}

func (d *Directory) LoadFiles(apps []string) ([]*File, error) {
	g, _ := errgroup.WithContext(context.Background())

	results := make([]*File, 0, len(apps))

	for _, app := range apps {
		app := app

		g.Go(func() error {
			fullPath := filepath.Join(d.path, fmt.Sprintf("%s.desktop", app))

			// ignore deleted files
			if !utils.Exists(fullPath) {
				d.logger.Warnf("file missed %s", fullPath)
				return nil
			}

			result, err := LoadFile(fullPath)

			if err == nil {
				results = append(results, result)
			} else {
				err = errors.New(err)
			}

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.New(err)
	}

	return results, nil
}

func (d *Directory) SaveFiles(files []*File) error {
	g, _ := errgroup.WithContext(context.Background())

	for _, file := range files {
		file := file

		g.Go(func() error {
			return file.Save()
		})
	}

	if err := g.Wait(); err != nil {
		return errors.New(err)
	}

	return nil
}
