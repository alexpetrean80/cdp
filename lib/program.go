package lib

import (
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func runCmd(executable string, args ...string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func OpenShell() error {
	return runCmd(os.Getenv("SHELL"))
}

func OpenEditor() error {
	editor := viper.GetString("editor")
	return runCmd(editor, ".")
}

func OpenGithubPage() error {
	return runCmd("gh", "browse")
}
