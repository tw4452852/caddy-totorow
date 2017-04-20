package totorow

import (
	"bytes"
	"encoding/xml"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tw4452852/storage"
	"golang.org/x/tools/blog/atom"
)

func init() {
	httpserver.RegisterDevDirective("totorow", "")

	caddy.RegisterPlugin("totorow", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	t := &Totorow{}

	err := parse(c, t)
	if err != nil {
		return err
	}

	registerTemplateFuncs()

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		t.Next = next
		return t
	})

	return initStorage(t.RepoConfig)
}

func parse(c *caddy.Controller, t *Totorow) error {
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 1:
			t.RepoConfig = args[0]
			t.BaseURL = "/"
		case 2:
			t.RepoConfig = args[0]
			t.BaseURL = args[1]
		default:
			return c.ArgErr()
		}
	}
	return nil
}

func initStorage(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	storage.Init(abs)
	return nil
}

func registerTemplateFuncs() {
	httpserver.TemplateFuncs["getFull"] = getFull
	httpserver.TemplateFuncs["getOne"] = getOne
	httpserver.TemplateFuncs["getTags"] = getTags
	httpserver.TemplateFuncs["getRss"] = getRss
	httpserver.TemplateFuncs["highlight"] = highlight
	httpserver.TemplateFuncs["filterTag"] = filterTag
	httpserver.TemplateFuncs["filterSearch"] = filterSearch
}

func getFull() []storage.Poster {
	result, err := storage.Get()
	if err != nil {
		log.Printf("getFull failed: %v\n", err)
		return nil
	}

	sort.Sort(result)
	return result.Content
}

type postKey string

func (pk postKey) Key() string {
	return string(pk)
}

func getOne(key string) storage.Poster {
	result, err := storage.Get(postKey(key))
	if err != nil {
		log.Printf("getOne failed: %v\n", err)
		return nil
	}
	return result.Content[0]
}

type TagEntry struct {
	Name  string
	Count int
}

type Tags []*TagEntry

func (ts Tags) add(name string) Tags {
	for _, t := range ts {
		if t.Name == name {
			t.Count++
			return ts
		}
	}
	ts = append(ts, &TagEntry{name, 1})
	return ts
}

func getTags(posts []storage.Poster) (ts Tags) {
	for _, p := range posts {
		for _, t := range p.Tags() {
			ts = ts.add(t)
		}
	}
	return ts
}

func highlight(search, input string) string {
	if search == "" {
		return input
	}

	index := strings.Index(input, search)
	if index == -1 {
		return input
	}

	return input[:index] + "<span class=highlight>" + search + "</span>" +
		input[index+len(search):]
}

func filterTag(tag string, posts []storage.Poster) []storage.Poster {
	return Filter(posts, CheckTags(tag))
}

func filterSearch(search string, posts []storage.Poster) []storage.Poster {
	return Filter(posts, CheckAll(search))
}

func getRss() string {
	r, err := storage.Get()
	if err != nil {
		return err.Error()
	}

	if len(r.Content) == 0 {
		return "refresh later"
	}

	//sorted by time
	sort.Sort(r)

	//Init a common infos
	feed := &atom.Feed{
		Title: "Tw's blog",
		Link:  []atom.Link{{Href: "/rss"}},
		ID:    "/rss",
		Author: &atom.Person{
			Name:  "Tw",
			Email: "tw19881113@gmail.com",
		},
		//use the newest post time
		Updated: atom.Time(r.Content[0].Date()),
	}

	//fill the entries
	feed.Entry = make([]*atom.Entry, len(r.Content))
	for i := range feed.Entry {
		p := r.Content[i]
		feed.Entry[i] = &atom.Entry{
			Title:   string(p.Title()),
			ID:      p.Key(),
			Updated: atom.Time(p.Date()),
			Link:    []atom.Link{{Href: "/posts/" + p.Key()}},
			Content: &atom.Text{Body: string(p.Content()), Type: "html"},
		}
	}

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")
	err = enc.Encode(feed)
	if err != nil {
		return err.Error()
	}

	return b.String()
}
