package cmd

import (
	"log"
	"os"

	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
)

var (
	last    bool
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
	rootCmd.PersistentFlags().BoolVarP(&last, "last", "l", false, "Change to the last project.")
}
