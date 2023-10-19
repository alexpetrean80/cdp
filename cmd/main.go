package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/alexpetrean80/cdp/config"
	"github.com/alexpetrean80/cdp/project"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	app := &cli.App{
		Name:  "cdp",
		Usage: "Move between projects seamlessly",
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	config, err := config.New()

	if err != nil {
		return err
	}

	ch := findProjects(config)
	projectPath, err := fzfProjectPath(ch)

	if err != nil {
	}

	err = os.Chdir(projectPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := spawnShell(); err != nil {
		log.Fatal(err)
	}
}

func spawnShell() error {
	shell := exec.Command(os.Getenv("SHELL"))
	shell.Stdin = os.Stdin
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	return shell.Run()
}

func findProjects(config *config.Config) <-chan string {
	g := new(errgroup.Group)
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

	return ch
}

func fzfProjectPath(ch <-chan string) (string, error) {
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
