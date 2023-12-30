package executable

import (
	"os"
	"os/exec"
)

type Program interface {
	Open() error
}

type program struct {
	executable string
	args       []string
}

func New(exec string, args ...string) Program {
	return &program{executable: exec, args: args}
}

func (p program) Open() error {
	cmd := exec.Command(p.executable, p.args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
