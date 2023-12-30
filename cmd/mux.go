/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"slices"

	"github.com/alexpetrean80/cdp/executable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	mux    executable.Program
	muxCmd = &cobra.Command{
		Use:   "mux",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			execName := viper.GetString("multiplexer")
			fmt.Println(execName)
			if i := slices.Index([]string{"tmux", "screen", "zellij"}, execName); i == -1 {
				return fmt.Errorf("%s is not a supported multiplexer. valid options are tmux, screen and zellij", execName)
			}

			mux = executable.New(execName, args...)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return mux.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(muxCmd)
	muxCmd.Flags().StringP("executable", "e", "", "program to be executed (one of: tmux, screen, zellij)")
	if err := viper.BindPFlag("multiplexer", muxCmd.Flags().Lookup("executable")); err != nil {
		log.Fatal(err.Error())
	}
}
