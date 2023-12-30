package lib
import (
	"sync"

	"github.com/ktr0731/go-fuzzyfinder"
)

func FuzzyFindProjectPath(ch <-chan string) (string, error) {
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
