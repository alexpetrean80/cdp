package lib

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func GetFullPathOfDirs() ([]string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	var res []string

	for _, dir := range viper.GetStringSlice("source.dirs") {
		res = append(res, fmt.Sprintf("%s/%s", homeDir, dir))
	}

	return res, nil
}

func ChangeDirectory(last bool) error {
	projectPath, err := GetProjectPath(last)
	if err != nil {
		log.Fatal(err.Error())
	}

	return os.Chdir(projectPath)
}
