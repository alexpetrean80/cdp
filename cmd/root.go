package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	last    bool
	cfgFile string

	rootCmd = &cobra.Command{
		Use: "cdp",
		//	Short: "A brief description of your application",
		//	Long: `A longer description that spans multiple lines and likely contains
		//
		// examples and usage of using your application. For example:
		//
		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := lib.ChangeDirectory(last); err != nil {
				log.Fatal(err.Error())
			}
		},
		// TODO find a way to specify default command in config
		// Run: func(cmd *cobra.Command, args []string) {
		//
		// },
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&last, "last", "l", false, "Change to the last project.")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "c", "config file (default is $HOME/.config/cdp/config.yaml")
}

func initConfig() {
	if cfgFile == "" {
		cfgFile = fmt.Sprintf("%s/.config/cdp/config.yaml", os.Getenv("HOME"))
	}

	viper.SetConfigFile(cfgFile)

}
