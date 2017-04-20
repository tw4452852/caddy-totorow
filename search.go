package totorow

import (
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/tw4452852/storage"
)

type Checker func(post storage.Poster) bool

func CheckAll(require string) Checker {
	return func(p storage.Poster) bool {
		if strings.Contains(string(p.Title()), require) {
			return true
		}
		if tagsContain(p.Tags(), require) {
			return true
		}
		if strings.Contains(string(p.Content()), require) {
			return true
		}
		return false
	}
}

func CheckTags(tag string) Checker {
	return func(p storage.Poster) bool {
		if tagsContain(p.Tags(), tag) {
			return true
		}
		return false
	}
}

func tagsContain(tags []string, require string) bool {
	for _, tag := range tags {
		if tag == require {
			return true
		}
	}
	return false
}

var workers = runtime.NumCPU()

// Filter return the Posters than contain `require`
// search sequence: title, tag, content
func Filter(all []storage.Poster, chk Checker) []storage.Poster {
	originCount := len(all)
	if originCount == 0 {
		return nil
	}

	in, out := make(chan storage.Poster, originCount), make(chan storage.Poster, originCount)
	// fill the input channel
	for _, p := range all {
		in <- p
	}
	close(in)

	// start workers
	waiter := new(sync.WaitGroup)
	for i := 0; i < workers; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			work(in, out, chk)
		}()
	}

	// collect the result
	waiter.Wait()
	resultCount := len(out)
	r := make([]storage.Poster, resultCount)
	for i := 0; i < resultCount; i++ {
		r[i] = <-out
	}

	ret := &storage.Result{r}
	sort.Sort(ret)
	return ret.Content
}

func work(in <-chan storage.Poster, out chan<- storage.Poster, chk Checker) {
	for p := range in {
		if chk(p) {
			out <- p
		}
	}
}
