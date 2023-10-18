package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
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

func main() {
	rootDir := "/Users/alexp/Repos"
	projects, err := findProjects(rootDir)
	if err != nil {
		log.Fatal(err)
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
