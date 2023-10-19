package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexpetrean80/cdp/config"
	"github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sync/errgroup"
)

func isHiddenDir(dir string) bool {
	return dir[0] == '.'
}

func findProjects(rootDir string, ch chan string, g *errgroup.Group) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if !entry.IsDir() {
			continue
		}

		if isHiddenDir(entryName) {
			if entryName == ".git" {
				ch <- rootDir
			}
		} else {
			g.Go(func() error {
				return findProjects(fmt.Sprintf("%s/%s", rootDir, entryName), ch, g)
			})
		}
	}
	return nil
}

func main() {
	g := new(errgroup.Group)

	config, err := config.New()

	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan string, 10)

	for _, dir := range config.Dirs() {
		g.Go(func(rootDir string) func() error {
			return func() error {
				return findProjects(rootDir, ch, g)
			}
		}(dir))
	}

	projects := []string{}

	go func() {
		defer close(ch)

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		proj, ok := <-ch
		if !ok {
			break
		}
		projects = append(projects, proj)
	}

	projId, err := fuzzyfinder.Find(
		projects, func(i int) string {
			return projects[i]
		})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(projects[projId])
}
