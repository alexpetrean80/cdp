/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		last, err := cmd.Flags().GetBool("last")
		if err != nil {
			log.Fatal(err.Error())
		}

		projectPath, err := lib.GetProjectPath(last)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = os.Chdir(projectPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := lib.SpawnProgram(os.Getenv("SHELL"), nil); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
