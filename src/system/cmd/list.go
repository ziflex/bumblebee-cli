package cmd

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/urfave/cli"
	"github.com/ziflex/bumblebee-gnome/src/core"
	"github.com/ziflex/bumblebee-gnome/src/core/fs"
	"github.com/ziflex/bumblebee-gnome/src/system/logging"
	"github.com/ziflex/bumblebee-gnome/src/system/storage"
	"github.com/ziflex/bumblebee-gnome/src/system/utils"
	"strings"
	"sort"
)

var (
	ERR_LIST_CMD = errors.New("failed to retrieve entries")
)

func NewListCommand(logger *logging.Logger, entryRepo storage.EntryRepository, settingsRepo storage.SettingsRepository) *cli.Command {
	return &cli.Command{
		Name:    "ls",
		Usage:   "show list of registered applications",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "a",
				Usage: "show all avialable applications",
			},
		},
		Action: func(ctx *cli.Context) error {
			entries, err := entryRepo.GetAll()

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_LIST_CMD
			}

			settings, err := settingsRepo.Get()

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_LIST_CMD
			}

			dir := fs.NewDirectory(logger, settings.Directory)
			all := ctx.NumFlags() > 0

			if !all {
				names := make([]string, 0, len(entries))

				for _, entry := range entries {
					names = append(names, entry.Name)
				}

				files, err := dir.LoadFiles(names)

				if err != nil {
					logger.Error(utils.ErrorStack(err))
					return ERR_LIST_CMD
				}

				filesMap := make(map[string]*fs.File)

				for _, file := range files {
					filesMap[file.Name()] = file
				}

				if err := printSimple(entries, filesMap); err != nil {
					logger.Error(utils.ErrorStack(err))
					return ERR_LIST_CMD
				}
			} else {
				names, err := dir.List()

				if err != nil {
					logger.Error(utils.ErrorStack(err))
					return ERR_LIST_CMD
				}

				printFull(entries, names)
			}

			return nil
		},
	}
}

func printSimple(entries map[string]*core.Entry, files map[string]*fs.File) error {
	list := make(map[string]string)
	names := make([]string, 0, len(entries))
	longestName := 0
	var err error

	for _, entry := range entries {
		file, ok := files[entry.Name]

		names = append(names, entry.Name)

		if longestName < len(entry.Name) {
			longestName = len(entry.Name)
		}

		if !ok {
			list[entry.Name] = "-"
			continue
		}

		values, failure := file.GetValues()

		if failure != nil {
			err = failure
			break
		}

		synced := true

		// check whether all values are synced
		for _, value := range values {
			if !core.IsGPUEnabled(value) {
				synced = false
				break
			}
		}

		list[entry.Name] = fmt.Sprintf("%t", synced)
	}

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("NAME %s SYNCED", getSpaces(longestName, "NAME")))

	sort.Strings(names)

	for _, name := range names {
		synced := list[name]

		fmt.Println(fmt.Sprintf("%s %s %s", name, getSpaces(longestName, name), synced))
	}

	return nil
}

func printFull(entries map[string]*core.Entry, names []string) {
	list := make(map[string]bool)
	longestName := 0

	// get formatting information and whether app is enabled to use GPU
	for _, name := range names {
		_, used := entries[name]

		if longestName < len(name) {
			longestName = len(name)
		}

		list[name] = used
	}

	sort.Strings(names)

	fmt.Println(fmt.Sprintf("NAME %s ENABLED", getSpaces(longestName, "NAME")))

	for _, name := range names {
		enabled := list[name]

		fmt.Println(fmt.Sprintf("%s %s %t", name, getSpaces(longestName, name), enabled))
	}
}

func getSpaces(max int, name string) string {
	if max == 0 {
		return strings.Repeat(" ", len(name))
	}

	return strings.Repeat(" ", max-len(name))
}
