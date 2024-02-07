package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexpetrean80/cdp/cmd"
	"github.com/spf13/viper"
)

func main() {
	cmd.Execute()
}

func init() {
	var configPath string

	if cp := os.Getenv("CDPCONFIG"); cp != "" {
		configPath = cp
	} else {
		configPath = fmt.Sprintf("%s/.config/cdp/config.yaml", os.Getenv("HOME"))
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
