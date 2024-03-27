package lib

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alexpetrean80/cdp/lib/finder"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// GetFullPathOfDirs is a function that returns the full path of the directories
// specified in the configuration file.
// It assumes that the directories are relative to the home directory.
func GetFullPathOfDirs() ([]string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	var res []string

	for _, dir := range viper.GetStringSlice("source.dirs") {
		res = append(res, fmt.Sprintf("%s/%s", homeDir, dir))
	}

	return res, nil
}

// ChangeDirectory is a function that changes the current working directory to the
// project directory specified by the name.
func ChangeDirectory(name string, last bool) error {
	projectPath, err := GetProjectPath(name, last)
	if err != nil {
		log.Fatal(err.Error())
	}

	return os.Chdir(projectPath)
}

// ReadLastProject is a function that reads the last project from the $HOME/.local/share/cdp_last file.
func ReadLastProject() (string, error) {
	file, err := os.Open(fmt.Sprintf("%s/.local/share/cdp_last", os.Getenv("HOME")))
	if err != nil {
		return "", err
	}

	proj, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(proj), nil
}

// WriteLastProject is a function that writes the last project to the $HOME/.local/share/cdp_last file.
func WriteLastProject(projectPath string) error {
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

// GetProjectPath is a function that returns the project path retrieved from finder.
func GetProjectPath(name string, last bool) (string, error) {
	if last {
		projectPath, err := ReadLastProject()
		log.Println(projectPath)
		if err == nil && projectPath != "" {
			return projectPath, nil
		}
	}

	ch := FindProjects(name)

	var projects []string

	for proj := range ch {
		projects = append(projects, proj)
	}

	var projectPath string
	switch len(projects) {
	case 0:
		return "", fmt.Errorf("No project found.")
	case 1:
		projectPath = projects[0]
	default:
		projId, err := fuzzyfinder.Find(
			projects,
			func(i int) string {
				return projects[i]
			})
		projectPath = projects[projId]
		if err != nil {
			return "", err
		}
	}

	if err := WriteLastProject(projectPath); err != nil {
		return projectPath, err
	}

	return projectPath, nil
}

// FindProjects is a function that searches for the projects with a filter.
// If the filter is empty, it returns all the projects found in the directories
// specified in the configuration file.
func FindProjects(name string) <-chan string {
	g := new(errgroup.Group)
	ch := make(chan string, 10)

	dirs, err := GetFullPathOfDirs()
	if err != nil {
		log.Fatal(err)
	}
	markers := viper.GetStringSlice("source.project_markers")
	for _, dir := range dirs {
		if !isDirectory(dir) {
			break
		}

		pf := finder.New(dir, markers, ch, g)
		g.Go(func(rootDir string) func() error {
			return func() error {
				return pf.Find(name)
			}
		}(dir))
	}

	go func() {
		defer close(ch)

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}
	}()

	return ch
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fileInfo.Mode().IsDir() {
		fmt.Printf("%s is not a directory, please check your config.\n", path)
		return false
	}

	return true
}
