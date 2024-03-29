package finder

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Finder struct {
	RootDir string
	Markers map[string]struct{}
	ResCh   chan string
	Group   *errgroup.Group
}

func New(rootDir string, markers []string, resCh chan string, g *errgroup.Group) *Finder {
	pf := Finder{}
	pf.RootDir = rootDir

	pf.Markers = make(map[string]struct{})
	for _, marker := range markers {
		pf.Markers[marker] = struct{}{}
	}

	pf.ResCh = resCh
	pf.Group = g

	return &pf
}

func (pf Finder) Find() error {
	return pf.findRec(pf.RootDir)
}

func (pf Finder) findRec(rootDir string) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := strings.Trim(entry.Name(), "/")

		if _, ok := pf.Markers[entryName]; ok {
			pf.ResCh <- rootDir
			return nil
		}

		if !isHiddenDir(entry) && entry.IsDir() {
			pf.Group.Go(func() error {
				return pf.findRec(fmt.Sprintf("%s/%s", rootDir, entryName))
			})
		}
	}

	return nil
}

func isHiddenDir(entry fs.DirEntry) bool {
	name := entry.Name()
	return name[0] == '.' && entry.IsDir()
}
