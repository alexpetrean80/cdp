package lib

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alexpetrean80/cdp/finder"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

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

func ChangeDirectory(name string, last bool) error {
	projectPath, err := GetProjectPath(name, last)
	if err != nil {
		log.Fatal(err.Error())
	}

	return os.Chdir(projectPath)
}

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
		fmt.Println("dir:", dir)
		fmt.Println("name:", name)

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
