/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/alexpetrean80/cdp/executable"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var (
	shellExecPath string
	shell         executable.Program
	shellCmd      = &cobra.Command{
		Use: "shell",
		// 	Short: "A brief description of your command",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:
		//
		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if shellExecPath == "" {
				shellExecPath = os.Getenv("SHELL")
			}

			if i := slices.Index([]string{"bash", "sh", "zsh", "fish"}); i == -1 {
				return fmt.Errorf("%s is not a supported shell. valid options are sh, bash, zsh and fish")
			}

			shell = executable.New(shellExecPath, args...)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return shell.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(shellCmd)
	muxCmd.Flags().StringVarP(&shellExecPath, "shell", "s", "", "shell to be opened (defaults to $SHELL)")
}
