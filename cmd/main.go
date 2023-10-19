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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "Open the project in the configured editor.",
			},
			&cli.BoolFlag{
				Name:    "browser",
				Aliases: []string{"b"},
				Usage:   "Open the project in the browser. (github-cli required)",
			},
		},
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
		return err
	}

	err = os.Chdir(projectPath)
	if err != nil {
		return err
	}

	if ctx.NumFlags() == 0 {
		if err := spawnProgram(os.Getenv("SHELL"), nil); err != nil {
			return err
		}
	}
	if ctx.Bool("open") {
		if err := spawnProgram(config.Editor, []string{"."}); err != nil {
			return err
		}
	}
	if ctx.Bool("browser") {
		fmt.Println("open in browser")
		if err := exec.Command("gh", "repo", "view", "--web").Run(); err != nil {
			return err
		}
	}

	return nil
}


func spawnProgram(executable string, args []string) error {
	cmd := exec.Command(executable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
