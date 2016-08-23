package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Unknwon/goconfig"
	"github.com/codegangsta/cli"
)

func IfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}

func FindSectionsByKey(file *goconfig.ConfigFile, name string) ([]string, error) {
	sections := file.GetSectionList()

	result := []string{}

	for _, sectionName := range sections {
		section, err := file.GetSection(sectionName)
		IfError(err)

		if section[name] != "" {
			result = append(result, sectionName)
		}
	}

	return result, nil
}

func CheckPermissions(path string) error {
	info, err := os.Lstat(path)

	if err != nil {
		return err
	}

	perm := info.Mode().Perm()

	fmt.Println("Permissions:", perm)

	if string(perm) != "-rwxr-xr-x" {
		return nil
		//return errors.New("No persmissions")
	}

	return nil
}

func Do(value string) string {
	if strings.HasPrefix(value, "primusrun") {
		return value
	}

	return "primusrun " + value
}

func Revert(value string) string {
	if !strings.HasPrefix(value, "primusrun") {
		return value
	}

	return strings.Replace(value, "primusrun", "", 1)
}

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "Update .desktop"
	app.Usage = "Updated .desktop files in order to run apps wiht 'primusrun' command"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "revert",
			Usage: "reverts operation",
		},
		cli.StringFlag{
			Name:  "pattern",
			Value: "",
			Usage: "files look up pattern",
		},
		cli.StringFlag{
			Name:  "files",
			Value: "",
			Usage: "coma separates list of target files",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("pattern") == "" && c.String("files") == "" {
			fmt.Println("Provide either look up pattern or list of files")
			os.Exit(0)
			return
		}

		path := "/usr/share/applications"

		err := CheckPermissions(path)
		IfError(err)

		files, err := ioutil.ReadDir(path)

		if err != nil {
			return
		}

		pattern := c.String("pattern")
		list := []string{}

		if c.String("files") != "" {
			content, err := ioutil.ReadFile(c.String("files"))

			IfError(err)

			if content == nil {
				IfError(errors.New("Passed file is empty"))
			}

			list = strings.Split(string(content), "\n")

		}

		processed := 0
		revert := c.Bool("revert")

		for _, file := range files {
			name := strings.TrimSpace(file.Name())
			matched := false
			if pattern != "" {
				matched, err = filepath.Match(pattern, name)
				IfError(err)
			} else {
				for _, f := range list {
					targetFileName := strings.TrimSpace(f)
					matched, err = filepath.Match(targetFileName, name)

					if matched == false {
						matched = targetFileName == name
					}

					if matched {
						break
					}
				}
			}

			if !matched {
				continue
			}

			fullName := filepath.Join(path, file.Name())
			cfg, err := goconfig.LoadConfigFile(fullName)
			IfError(err)

			sections, err := FindSectionsByKey(cfg, "Exec")

			IfError(err)

			toSave := false

			for _, section := range sections {
				value, err := cfg.GetValue(section, "Exec")

				IfError(err)

				newValue := value
				if revert == false {
					newValue = Do(value)
				} else {
					newValue = Revert(value)
				}

				if newValue == value {
					continue
				}

				toSave = true
				cfg.SetValue(section, "Exec", newValue)
			}

			if toSave {
				processed += 1
				fmt.Println(fullName)
				goconfig.SaveConfigFile(cfg, fullName)
			}
		}

		fmt.Println("Processed " + strconv.Itoa(processed) + " file(s)")
	}

	app.Run(os.Args)
}
