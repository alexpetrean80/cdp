/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			if shellExecPath == "" {
				s := strings.Split(os.Getenv("SHELL"), "/")
				slices.Reverse(s)
				shellExecPath = s[0]
			}

			if i := slices.Index([]string{"bash", "sh", "zsh", "fish"}, shellExecPath); i == -1 {
				return fmt.Errorf("%s is not a supported shell. valid options are sh, bash, zsh and fish", shellExecPath)
			}

			shell = executable.New(shellExecPath, args...)

			return shell.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&shellExecPath, "shell", "s", "", "shell to be opened (defaults to $SHELL)")
}
