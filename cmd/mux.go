/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/alexpetrean80/cdp/lib"
	"github.com/spf13/cobra"
)

var muxCmd = &cobra.Command{
	Use:   "mux",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lib.OpenMultiplexer()
	},
}

func init() {
	rootCmd.AddCommand(muxCmd)
}
