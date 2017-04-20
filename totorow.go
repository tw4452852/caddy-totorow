package totorow

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Totorow struct {
	Next       httpserver.Handler
	BaseURL    string
	RepoConfig string
}

func (t *Totorow) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if httpserver.Path(r.URL.Path).Matches(t.BaseURL) {
		// handle post's image
		if httpserver.Path(r.URL.Path).Matches(t.BaseURL + "images/") {
			pathBegin := strings.LastIndex(r.URL.Path, "/")
			keyBegin := strings.LastIndex(r.URL.Path[:pathBegin], "/")
			if pathBegin == -1 || keyBegin == -1 {
				return http.StatusInternalServerError, errors.New("image url invalid")
			}
			key := r.URL.Path[keyBegin+1 : pathBegin]
			path := r.URL.Path[pathBegin+1:]
			if path == "" || key == "" {
				return http.StatusInternalServerError, errors.New("image path or key invalid")
			}
			return serveImage(w, r, path, key)
		}
	}

	return t.Next.ServeHTTP(w, r)
}

func serveImage(w http.ResponseWriter, r *http.Request, path, key string) (int, error) {
	p := getOne(key)
	if p == nil {
		return http.StatusInternalServerError, errors.New("can't find related post")
	}
	sta := p.Static(path)
	defer sta.Close()

	io.Copy(w, sta)
	return 0, nil
}
