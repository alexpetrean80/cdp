package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/alexpetrean80/cdp/lib/executable"
	"github.com/spf13/cobra"
)

var (
	shellExecPath string
	shell         executable.Program
	shellCmd      = &cobra.Command{
		Use:   "shell",
		Short: "Open a shell in the project's directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if shellExecPath == "" {
				s := strings.Split(os.Getenv("SHELL"), "/")
				slices.Reverse(s)
				shellExecPath = s[0]
			}

			if i := slices.Index([]string{"bash", "sh", "zsh", "fish"}, shellExecPath); i == -1 {
				return fmt.Errorf(
					"%s is not a supported shell. valid options are sh, bash, zsh and fish",
					shellExecPath,
				)
			}

			shell = executable.New(shellExecPath, args...)

			return shell.Open()
		},
	}
)

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().
		StringVarP(&shellExecPath, "shell", "s", "", "shell to be opened (defaults to $SHELL)")
}
