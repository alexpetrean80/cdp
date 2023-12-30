package cmd

import (
	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
)

// browserCmd represents the browser command
var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lib.OpenGithubPage()
	},
}

func init() {
	rootCmd.AddCommand(browserCmd)
}
