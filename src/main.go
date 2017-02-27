package main

import (
	"fmt"
	"github.com/ziflex/bumblebee-cli/src/system"
	"os"
)

func main() {

	if os.Geteuid() != 0 {
		handleError(os.ErrPermission)
		return
	}

	app, err := system.NewApplication()

	if err != nil {
		handleError(err)
		return
	}

	err = app.Run(os.Args)

	if err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	printError(err)
	os.Exit(1)
}

func printError(err error) {
	fmt.Println(fmt.Sprintf("Error: %s", err.Error()))
}
