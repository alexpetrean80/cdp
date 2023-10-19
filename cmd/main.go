package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/alexpetrean80/cdp/config"
	"github.com/alexpetrean80/cdp/project"
	"github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sync/errgroup"
)

func main() {
	g := new(errgroup.Group)

	config, err := config.New()

	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan string, 10)

	for _, dir := range config.Dirs() {
		pf := project.New(dir, config.Markers(), ch, g)
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

	mtx := new(sync.Mutex)
	projects := []string{}
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

	projectPath, err := fzfProjectPath(&projects, mtx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(projectPath)
}

func fzfProjectPath(projects *[]string, mtx sync.Locker) (string, error) {
	projId, err := fuzzyfinder.Find(
		projects, func(i int) string {
			return (*projects)[i]
		},
		fuzzyfinder.WithHotReloadLock(mtx),
	)

	if err != nil {
		return "", err
	}

	return (*projects)[projId], nil
}
