package finder

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"golang.org/x/sync/errgroup"
)

// Finder is a struct that contains the root directory to start the search from,
// the markers to search for, the result channel to send the results to, and the
// errgroup to manage the goroutines.
type Finder struct {
	RootDir string
	Markers map[string]struct{}
	ResCh   chan string
	Group   *errgroup.Group
}

// New creates a new Finder struct with the given root directory, markers, result
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

// Find starts the search from the root directory and sends the results to the result channel.
func (pf Finder) Find(name string) error {
	return pf.findRec(pf.RootDir, name)
}

// findRec is a recursive function that searches for the given name in the root directory.
// If the name is found, it sends the result to the result channel.
// If the entry is a directory, it starts a new goroutine to search in the directory.
// If the entry is a hidden directory, it skips the directory.
func (pf Finder) findRec(rootDir, name string) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryName := strings.Trim(entry.Name(), "/")

		if _, ok := pf.Markers[entryName]; ok {
			if fuzzy.MatchNormalized(name, rootDir) {
				pf.ResCh <- rootDir
				return nil
			}
		}

		if !isHiddenDir(entry) && entry.IsDir() {
			pf.Group.Go(func() error {
				return pf.findRec(fmt.Sprintf("%s/%s", rootDir, entryName), name)
			})
		}
	}

	return nil
}

func isHiddenDir(entry fs.DirEntry) bool {
	name := entry.Name()
	return name[0] == '.' && entry.IsDir()
}
