package utils

import (
	"github.com/go-errors/errors"
)

func ErrorStack(err error) string {
	if err != nil {
		if e, ok := err.(*errors.Error); ok {
			return e.ErrorStack()
		} else {
			return errors.New(err).ErrorStack()
		}
	}

	return ""
}

func ErrorStackf(message string, values ...interface{}) string {
	return ErrorStack(errors.Errorf(message, values...))
}
