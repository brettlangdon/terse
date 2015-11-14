package terse

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/groupcache/lru"
)

type Handler struct {
	cache *lru.Cache
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler.HandleGet(w, r)
	case "POST":
		handler.HandlePost(w, r)
	default:
		http.Error(w, fmt.Sprintf("Method \"%s\" Not Allowed", r.Method), http.StatusMethodNotAllowed)
	}
}

func (handler *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	url := strings.Trim(r.URL.Path, "/")
	code := strings.SplitN(url, "/", 2)[0]
	value, ok := handler.cache.Get(code)
	if ok {
		switch url := value.(type) {
		case string:
			http.Redirect(w, r, url, http.StatusMovedPermanently)
			return
		}
	}
	http.NotFound(w, r)
}

func (handler *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	reader := bufio.NewScanner(r.Body)
	reader.Scan()
	url, err := url.ParseRequestURI(reader.Text())
	if err != nil {
		panic(err)
	}
	code := GetShortCode([]byte(url.String()))
	handler.cache.Add(code, url.String())
	fmt.Fprintf(w, code)
}

func NewServer(bind string, maxEntries int) *http.Server {
	return &http.Server{
		Addr: bind,
		Handler: &Handler{
			cache: lru.New(maxEntries),
		},
	}
}
