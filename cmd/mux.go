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
		Short: "Open project in a mux session",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			muxExecPath := viper.GetString("multiplexer")
			if i := slices.Index([]string{"tmux", "screen", "zellij"}, muxExecPath); i == -1 {
				return fmt.Errorf("%s is not a supported multiplexer. valid options are tmux, screen and zellij", muxExecPath)
			}

			mux = executable.New(muxExecPath, args...)
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
