package cmd

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/urfave/cli"
	"github.com/ziflex/bumblebee-cli/src/core"
	"github.com/ziflex/bumblebee-cli/src/core/fs"
	"github.com/ziflex/bumblebee-cli/src/system/logging"
	"github.com/ziflex/bumblebee-cli/src/system/storage"
	"github.com/ziflex/bumblebee-cli/src/system/utils"
)

var (
	ERR_UPDATE_CMD = errors.New("failed to update")
)

func NewSyncCommand(logger *logging.Logger, entryRepo storage.EntryRepository, settingsRepo storage.SettingsRepository) *cli.Command {
	return &cli.Command{
		Name:    "sync",
		Usage:   "update files of registered applications",
		Aliases: []string{"s"},
		Action: func(ctx *cli.Context) error {
			entries, err := entryRepo.GetAll()

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_UPDATE_CMD
			}

			if len(entries) == 0 {
				fmt.Println("No registered applications")
				return nil
			}

			settings, err := settingsRepo.Get()

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_UPDATE_CMD
			}

			list := make([]*core.Entry, 0, len(entries))

			for _, entry := range entries {
				list = append(list, entry)
			}

			transformer := core.NewTransformer(logger, fs.NewDirectory(logger, settings.Directory))

			results, err := transformer.Do(list, settings.Prefix, false)

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_UPDATE_CMD
			}

			if len(results) == 0 {
				fmt.Println("All files are up to date")
				return nil
			}

			fmt.Println("Synced:")

			reported := make(map[string]bool)

			for _, name := range results {
				if _, ok := reported[name]; !ok {
					fmt.Println(name)
					reported[name] = true
				}
			}

			return nil
		},
	}
}
