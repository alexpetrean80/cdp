package lib

import (
	"fmt"
	"os"
)

func writeLastProjectToFile(projectPath string) error {
	file, err := os.Create(fmt.Sprintf("%s/.local/share/cdp_last", os.Getenv("HOME")))
	if err != nil {
		return err
	}

	_, err = file.WriteString(projectPath)
	if err != nil {
		return err
	}

	return nil
}
