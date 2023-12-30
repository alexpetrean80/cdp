package lib

import (
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func spawnProgram(executable string, args []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SpawnShell() error {
	return spawnProgram(os.Getenv("SHELL"), nil)
}

func SpawnEditor() error {
	editor := viper.GetString("editor")

	return spawnProgram(editor, []string{"."})

}

func SpawnMux() error {
	mux := viper.GetString("multiplexer")
	return spawnProgram(mux, nil)
}
