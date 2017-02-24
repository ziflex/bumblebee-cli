package system

import (
	"database/sql"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/urfave/cli"
	"github.com/ziflex/bumblebee-gnome/src/system/cmd"
	"github.com/ziflex/bumblebee-gnome/src/system/initialization"
	"github.com/ziflex/bumblebee-gnome/src/system/initialization/initializers"
	"github.com/ziflex/bumblebee-gnome/src/system/logging"
	"github.com/ziflex/bumblebee-gnome/src/system/storage"
	"github.com/ziflex/bumblebee-gnome/src/system/storage/sqlite"
	"github.com/ziflex/bumblebee-gnome/src/system/utils"
	"path"
	"strings"
)

type Application struct {
	engine       *cli.App
	db           *sql.DB
	initManager  *initialization.InitManager
	initializers map[string]initialization.Initializer
	commands     []cli.Command
}

func NewApplication() (*Application, error) {
	var err error
	app := &Application{}

	app.engine = cli.NewApp()
	app.engine.Version = "2.0.0"
	app.engine.Name = "bumblebee-gnome"
	app.engine.Usage = "Manager for bumblebee dependant applications"

	logsDir := fmt.Sprintf("/var/log/%s/", strings.ToLower(app.engine.Name))

	if err = utils.EnsureDirectory(logsDir); err != nil {
		return nil, err
	}

	logger := logging.NewLogger(&lumberjack.Logger{
		Dir:        logsDir,
		MaxSize:    50 * lumberjack.Megabyte, // megabytes
		MaxBackups: 2,
		MaxAge:     28, //days
	})

	dbDir := fmt.Sprintf("/var/lib/%s/", strings.ToLower(app.engine.Name))

	if err = utils.EnsureDirectory(dbDir); err != nil {
		return nil, err
	}

	db, err := storage.OpenDb(path.Join(dbDir, "database.db"))

	if err != nil {
		logger.Fatalf("Failed to open db: %s", err.Error())
		return nil, err
	}

	app.db = db

	entries := sqlite.NewEntryRepository(storage.ENTRY_TABLE, db)
	settings := sqlite.NewSettingsRepository(storage.SETTINGS_TABLE, db)

	app.commands = []cli.Command{
		*cmd.NewListCommand(logger, entries, settings),
		*cmd.NewAddCommand(logger, entries, settings),
		*cmd.NewRemoveCommand(logger, entries, settings),
		*cmd.NewSyncCommand(logger, entries, settings),
	}

	app.initManager = initialization.NewInitManager(logger)
	app.initializers = map[string]initialization.Initializer{
		"database": initializers.NewDatabaseInitializer(logger, app.db),
	}

	return app, err
}

func (app *Application) Run(arguments []string) error {
	var err error

	defer app.db.Close()

	if err = app.initManager.Run(app.initializers); err != nil {
		return err
	}

	app.engine.Commands = app.commands

	return app.engine.Run(arguments)
}
