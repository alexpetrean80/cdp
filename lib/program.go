package lib

import (
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func runCmd(executable string, args ...string) error {
}

type Program interface {
	Open() error
}

type program struct {
	executable string
	args       []string
}

func (p program) Open() error {
	cmd := exec.Command(p.executable, p.args...)
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

func OpenMultiplexer() error {
	mux := viper.GetString("multiplexer")
	return runCmd(mux)
}

func OpenGithubPage() error {
	return runCmd("gh", "browse")
}
