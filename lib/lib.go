package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/alexpetrean80/cdp/project"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// func getProjectName(projectPath string) string {
// 	split := strings.Split(projectPath, "/")
// 	return split[len(split)-1]
// }

func GetProjectPath(last bool) (string, error) {
  log.Println(last)
	if last {
		projectPath, err := TryGetLatestProject()
		log.Println(projectPath)
		if err == nil && projectPath != "" {
			return projectPath, nil
		}
	}

	ch := FindProjects()
	projectPath, err := FuzzyFindProjectPath(ch)
	if err != nil {
		return "", err
	}

	if err := writeLastProjectToFile(projectPath); err != nil {
		return projectPath, err
	}

	return projectPath, nil
}

func FindProjects() <-chan string {
	g := new(errgroup.Group)
	ch := make(chan string, 10)

	dirs, err := GetDirs()
	if err != nil {
		log.Fatal(err)
	}
	markers := viper.GetStringSlice("source.project_markers")
	for _, dir := range dirs {
		pf := project.NewFinder(dir, markers, ch, g)
		g.Go(func(rootDir string) func() error {
			return func() error {
				return pf.Find()
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

func FuzzyFindProjectPath(ch <-chan string) (string, error) {
	mtx := new(sync.Mutex)
	var projects []string
	go func() {
		for {
			proj, ok := <-ch
			if !ok {
				break
			}
			mtx.Lock()
			projects = append(projects, proj)
			mtx.Unlock()
		}
	}()

	projId, err := fuzzyfinder.Find(
		&projects, func(i int) string {
			return (projects)[i]
		},
		fuzzyfinder.WithHotReloadLock(mtx),
	)

	if err != nil {
		return "", err
	}

	return (projects)[projId], nil
}

func TryGetLatestProject() (string, error) {
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

func GetDirs() ([]string, error) {
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

func ChangeDirectory(last bool) error {
	projectPath, err := GetProjectPath(last)
	if err != nil {
		log.Fatal(err.Error())
	}

	return os.Chdir(projectPath)
}
