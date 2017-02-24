package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var DefaultOutput = os.Stdout

type (
	Logger struct {
		info  *log.Logger
		warn  *log.Logger
		error *log.Logger
		fatal *log.Logger
	}
)

func NewDefaultLogger() *Logger {
	return NewLogger(DefaultOutput)
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{
		info:  log.New(out, "[INFO] ", log.Ldate|log.Ltime),
		warn:  log.New(out, "[WARN] ", log.Ldate|log.Ltime),
		error: log.New(out, "[ERROR] ", log.Ldate|log.Ltime),
		fatal: log.New(out, "[FATAL] ", log.Ldate|log.Ltime),
	}
}

func (logger *Logger) Info(message string) {
	logger.info.Println(message)
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.info.Println(fmt.Sprintf(format, v...))
}

func (logger *Logger) Warn(message string) {
	logger.warn.Println(message)
}

func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.warn.Println(fmt.Sprintf(format, v...))
}

func (logger *Logger) Error(message string) {
	logger.error.Println(message)
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.error.Println(fmt.Sprintf(format, v...))
}

func (logger *Logger) Fatal(message string) {
	logger.fatal.Println(message)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.fatal.Println(fmt.Sprintf(format, v...))
}
