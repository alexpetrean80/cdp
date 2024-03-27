package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	last    bool
	cfgFile string
	name    string

	rootCmd = &cobra.Command{
		Use: "cdp",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if slices.Contains([]string{"last", "completion"}, cmd.Use) {
				return nil
			}
			return lib.ChangeDirectory(name, last)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	if cfgFile == "" {
		if cp := os.Getenv("CDPCONFIG"); cp != "" {
			cfgFile = cp
		} else {
			cfgFile = fmt.Sprintf("%s/.config/cdp/config.yaml", os.Getenv("HOME"))
		}
	}

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}

	last = false
	rootCmd.PersistentFlags().BoolVarP(&last, "last", "l", false, "Change to the last project.")

	rootCmd.PersistentFlags().
		StringVarP(&name, "name", "n", "", "Specify a substring from a project name to be used in the search. Filtration is done fuzzily. If only one result is returned, it will bypass the fuzzyfinder entirely.")

	rootCmd.PersistentFlags().
		StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/cdp/config.yaml")
}
