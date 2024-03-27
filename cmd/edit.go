package cmd

import (
	"log"

	"github.com/alexpetrean80/cdp/lib/executable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	editor executable.Program

	editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Open project in the editor.",
		RunE: func(cmd *cobra.Command, args []string) error {
			editorExecPath := viper.GetString("editor")
			editor = executable.New(editorExecPath, append(args, ".")...)
			return editor.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("executable", "e", "", "program to be used")
	if err := viper.BindPFlag("editor", editCmd.Flags().Lookup("executable")); err != nil {
		log.Fatal(err.Error())
	}
}
