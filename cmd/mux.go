package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/alexpetrean80/cdp/lib/executable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getSessionName() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	dir := strings.Split(pwd, "/")
	slices.Reverse(dir)

	return dir[0]
}

var (
	mux    executable.Program
	muxCmd = &cobra.Command{
		Use:   "mux",
		Short: "Open project in a mux session",
		RunE: func(cmd *cobra.Command, args []string) error {
			muxExecPath := viper.GetString("multiplexer")
			if i := slices.Index([]string{"tmux", "screen", "zellij"}, muxExecPath); i == -1 {
				return fmt.Errorf(
					"%s is not a supported multiplexer. valid options are tmux, screen and zellij, muxExecPath",
					muxExecPath,
				)
			}

			sessionName := getSessionName()
			args = append(args, "new")
			args = append(args, fmt.Sprintf("-s %s", sessionName))
			mux = executable.New(muxExecPath, args...)
			return mux.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(muxCmd)
	muxCmd.Flags().
		StringP("executable", "e", "", "program to be executed (one of: tmux, screen, zellij)")
	if err := viper.BindPFlag("multiplexer", muxCmd.Flags().Lookup("executable")); err != nil {
		log.Fatal(err.Error())
	}
}
