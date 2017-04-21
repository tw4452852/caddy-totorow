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
			key, path, err := getKeyAndPath(strings.TrimPrefix(r.URL.Path, t.BaseURL+"images/"))
			if err != nil {
				return http.StatusInternalServerError, err
			}
			return serveImage(w, r, key, path)
		}
	}

	return t.Next.ServeHTTP(w, r)
}

var (
	imageURLInvalid = errors.New("image url invalid")
)

func getKeyAndPath(url string) (string, string, error) {
	pathBegin := strings.Index(url, "/")
	if pathBegin == -1 {
		return "", "", imageURLInvalid
	}
	key := url[:pathBegin]
	path := url[pathBegin+1:]
	if path == "" || key == "" {
		return "", "", imageURLInvalid
	}
	return key, path, nil
}

func serveImage(w http.ResponseWriter, r *http.Request, key, path string) (int, error) {
	p := getOne(key)
	if p == nil {
		return http.StatusInternalServerError, errors.New("can't find related post")
	}
	sta := p.Static(path)
	defer sta.Close()

	_, err := io.Copy(w, sta)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
