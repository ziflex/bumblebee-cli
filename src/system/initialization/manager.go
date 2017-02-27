package initialization

import (
	"github.com/go-errors/errors"
	"github.com/ziflex/bumblebee-ui/src/system/logging"
	"github.com/ziflex/bumblebee-ui/src/system/utils"
)

type (
	Initializer interface {
		Run() error
	}

	InitManager struct {
		logger *logging.Logger
	}
)

func NewInitManager(logger *logging.Logger) *InitManager {
	return &InitManager{
		logger: logger,
	}
}

func (manager *InitManager) Run(initializers map[string]Initializer) error {
	var err error

	for name, init := range initializers {
		initError := init.Run()

		if err != nil {
			err = errors.Errorf("Error occured during %s initializer: %s", name, initError.Error())
			manager.logger.Error(utils.ErrorStack(err))
			break
		}
	}

	return err
}
