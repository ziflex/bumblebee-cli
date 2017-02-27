package cmd

import (
	"bufio"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/urfave/cli"
	"github.com/ziflex/bumblebee-cli/src/core"
	"github.com/ziflex/bumblebee-cli/src/core/fs"
	"github.com/ziflex/bumblebee-cli/src/system/logging"
	"github.com/ziflex/bumblebee-cli/src/system/storage"
	"github.com/ziflex/bumblebee-cli/src/system/utils"
	"os"
	"strings"
)

var (
	ERR_SETTINGS_CMD   = errors.New("failed to set new setting")
	ERR_MISSED_VALUE   = errors.New("missed value")
	ERR_MISSED_KEY     = errors.New("missed key")
	ERR_MANY_ARGS      = errors.New("too many arguments")
	ERR_INVALID_PREFIX = errors.Errorf("prefix must be one of '%s' / '%s'", core.PRIMUSRUN, core.OPTIRUN)
	ERR_UPDATE_APPS    = errors.New("failed to update applications")
)

func NewSettingsCommand(logger *logging.Logger, entryRepo storage.EntryRepository, settingsRepo storage.SettingsRepository) *cli.Command {
	return &cli.Command{
		Name:    "setting",
		Usage:   "Manages bumblebee-ui settings",
		Aliases: []string{"s"},
		Subcommands: []cli.Command{
			{
				Name:  "set",
				Usage: "Sets a new setting value [key] [value]",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return ERR_MISSED_ARGS
					}

					if ctx.NArg() == 1 {
						return ERR_MISSED_VALUE
					}

					if ctx.NArg() > 2 {
						return ERR_MANY_ARGS
					}

					args := ctx.Args()

					key := strings.TrimSpace(args.Get(0))
					value := strings.TrimSpace(args.Get(1))

					switch key {
					case "prefix":
						return setPrefix(logger, entryRepo, settingsRepo, value)
					default:
						return errors.Errorf("invalid key '%s'", key)
					}
				},
			},
			{
				Name:  "get",
				Usage: "Gets a current seting value [key]",
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return ERR_MISSED_KEY
					}

					if ctx.NArg() > 1 {
						return ERR_MANY_ARGS
					}

					args := ctx.Args()

					key := strings.TrimSpace(args.Get(0))

					switch key {
					case "prefix":
						return getPrefix(logger, settingsRepo)
					default:
						return errors.Errorf("invalid key '%s'", key)
					}
				},
			},
		},
	}
}

func getPrefix(logger *logging.Logger, settingsRepo storage.SettingsRepository) error {
	settings, err := settingsRepo.Get()

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return ERR_SETTINGS_CMD
	}

	fmt.Println(settings.Prefix)

	return nil
}

func setPrefix(logger *logging.Logger, entryRepo storage.EntryRepository, settingsRepo storage.SettingsRepository, value string) error {
	if value != core.OPTIRUN && value != core.PRIMUSRUN {
		return ERR_INVALID_PREFIX
	}

	settings, err := settingsRepo.Get()

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return ERR_SETTINGS_CMD
	}

	if settings.Prefix == value {
		return nil
	}

	settings.Prefix = value

	err = settingsRepo.Set(settings)

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return ERR_SETTINGS_CMD
	}

	answer, err := ask("Do you want to update existing applications?", []string{"y", "n"})

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return errors.New("unexpected error")
	}

	if answer == "n" {
		return nil
	}

	entries, err := entryRepo.GetAll()

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return ERR_GET_ENTRIES
	}

	transformer := core.NewTransformer(logger, fs.NewDirectory(logger, settings.Directory))

	_, err = transformer.Do(core.EntryMapToList(entries), settings.Prefix, true)

	if err != nil {
		logger.Error(utils.ErrorStack(err))
		return ERR_UPDATE_APPS
	}

	return nil
}

func ask(question string, options []string) (string, error) {
	optionsStr := strings.Join(options, "/")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(fmt.Sprintf("%s (%s):", question, optionsStr))

	answer, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(answer)

	if !contains(options, result) {
		fmt.Println(fmt.Sprintf("Answer msut be one of (%s)", optionsStr))
		return ask(question, options)
	}

	return result, nil
}

func contains(list []string, targetValue string) bool {
	for _, currentValue := range list {
		if currentValue == targetValue {
			return true
		}
	}
	return false
}
