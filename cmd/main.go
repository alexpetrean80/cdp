package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sync/errgroup"
)

func isHiddenDir(dir string) bool {
	return dir[0] == '.'
}

func findProjects(rootDir string) ([]string, error) {
	projects := []string{}
	entries, err := os.ReadDir(rootDir)

	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryName := entry.Name()

		if !entry.IsDir() {
			continue
		}

		if isHiddenDir(entryName) {
			if entryName == ".git" {
				projects = append(projects, rootDir)
				return projects, nil
			} else {
				continue
			}
		}

		subProjects, err := findProjects(fmt.Sprintf("%s/%s", rootDir, entryName))
		if err != nil {
			return nil, err
		}
		projects = append(projects, subProjects...)
	}
	return projects, err
}

func findProjectsPar(rootDir string, ch chan string, g *errgroup.Group) error {
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
				return findProjectsPar(fmt.Sprintf("%s/%s", rootDir, entryName), ch, g)
			})
		}
	}
	return nil
}

func main() {
	g := new(errgroup.Group)

	rootDir := "/Users/alexp/Repos"
	// projects, err := findProjects(rootDir)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	ch := make(chan string, 10)

	g.Go(func() error {
		return findProjectsPar(rootDir, ch, g)
	})

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
			break // channel closed
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
