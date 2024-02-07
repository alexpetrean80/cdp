/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/alexpetrean80/cdp/executable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// editCmd represents the edit command
var (
	editor executable.Program

	editCmd = &cobra.Command{
		Use: "edit",
		// 	Short: "A brief description of your command",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:
		//
		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,
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
