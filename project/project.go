package project

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
)

type ProjectFinder struct {
	RootDir string
	Markers map[string]bool
	ResCh   chan string
	Group   *errgroup.Group
}

func New(rootDir string, markers []string, resCh chan string, g *errgroup.Group) *ProjectFinder {
	pf := ProjectFinder{}
	pf.RootDir = rootDir

	pf.Markers = make(map[string]bool)
	for _, marker := range markers {
		pf.Markers[marker] = true
	}

	pf.ResCh = resCh
	pf.Group = g

	return &pf
}

func (pf ProjectFinder) Find() error {
	return pf.findRec(pf.RootDir)
}

func (pf ProjectFinder) findRec(rootDir string) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := strings.Trim(entry.Name(), "/")

		if pf.Markers[entryName] {
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
