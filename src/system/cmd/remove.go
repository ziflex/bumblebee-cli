package cmd

import (
	"github.com/go-errors/errors"
	"github.com/urfave/cli"
	"github.com/ziflex/bumblebee-ui/src/core"
	"github.com/ziflex/bumblebee-ui/src/core/fs"
	"github.com/ziflex/bumblebee-ui/src/system/logging"
	"github.com/ziflex/bumblebee-ui/src/system/storage"
	"github.com/ziflex/bumblebee-ui/src/system/utils"
)

var (
	ERR_REMOVE_CMD = errors.New("failed to remove entry")
)

func NewRemoveCommand(logger *logging.Logger, entryRepo storage.EntryRepository, settingsRepo storage.SettingsRepository) *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Usage:   "Removes an application from the app registry",
		Aliases: []string{"r"},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return ERR_MISSED_ARGS
			}

			entries := make([]*core.Entry, 0, ctx.NArg())

			list := ctx.Args()

			for _, name := range ctx.Args() {
				entries = append(entries, core.NewEntry(name))
			}

			settings, err := settingsRepo.Get()

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_REMOVE_CMD
			}

			transformer := core.NewTransformer(logger, fs.NewDirectory(logger, settings.Directory))
			_, err = transformer.Revert(entries)

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_REMOVE_CMD
			}

			err = entryRepo.DeleteManyByName(list)

			if err != nil {
				logger.Error(utils.ErrorStack(err))
				return ERR_REMOVE_CMD
			}

			return nil
		},
	}
}
