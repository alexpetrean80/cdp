package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/alexpetrean80/cdp/config"
	"github.com/alexpetrean80/cdp/project"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	app := &cli.App{
		Name:                   "cdp",
		UseShortOptionHandling: true,
		Usage:                  "Move between projects seamlessly",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Open the project in the configured editor. (mutually exclusive with -t)",
			},
			&cli.BoolFlag{
				Name:    "browser",
				Aliases: []string{"b"},
				Usage:   "Open the project in the browser. github-cli required",
			},
			&cli.BoolFlag{
				Name:    "latest",
				Aliases: []string{"l"},
				Usage:   "Open the latest project",
			},
			&cli.BoolFlag{
				Name:    "tmux",
				Aliases: []string{"t"},
				Usage:   "Open the project in a new tmux session. (mutually exclusive with -o) tmux required.",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	c, err := config.New()

	if err != nil {
		return err
	}

	projectPath, err := getProjectPath(ctx, c)
	if err != nil {
		return err
	}

	err = os.Chdir(projectPath)
	if err != nil {
		return err
	}

	if ctx.NumFlags() == 0 || (ctx.NumFlags() == 2 && ctx.Bool("latest")) {
		if err = spawnProgram(os.Getenv("SHELL"), nil); err != nil {
			return err
		}
	}
	if ctx.Bool("open") {
		if err = spawnProgram(c.Editor, []string{"."}); err != nil {
			return err
		}
	} else if ctx.Bool("tmux") {
		if err = spawnProgram("tmux", []string{"new", "-s", getProjectName(projectPath)}); err != nil {
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

func getProjectName(projectPath string) string {
	split := strings.Split(projectPath, "/")
	return split[len(split)-1]
}
func getProjectPath(ctx *cli.Context, config *config.Config) (string, error) {
	if ctx.Bool("latest") {
		projectPath, err := tryGetLatestProject()
		if err == nil && projectPath != "" {
			return projectPath, nil
		}
	}

	ch := findProjects(config)
	projectPath, err := fzfProjectPath(ch)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fmt.Sprintf("%s/.config/cdp/latest", os.Getenv("HOME")))
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(projectPath)
	if err != nil {
		return "", err
	}

	return projectPath, nil
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

func tryGetLatestProject() (string, error) {
	file, err := os.Open(fmt.Sprintf("%s/.config/cdp/latest", os.Getenv("HOME")))
	if err != nil {
		return "", err
	}

	proj, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(proj), nil
}
