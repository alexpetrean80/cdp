package lib

import (
	"log"

	"github.com/alexpetrean80/cdp/project"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func GetProjectPath(last bool) (string, error) {
	log.Println(last)
	if last {
		projectPath, err := readLastProject()
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
