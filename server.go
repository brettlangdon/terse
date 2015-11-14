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
	// If they gave us something like `/<shortcode>/<who knows>`
	// then jst grab the first part
	code := strings.SplitN(url, "/", 2)[0]
	value, ok := handler.cache.Get(code)

	// If the url exists in the cache and is a string, redirect to it
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
	// Parse url from body
	reader := bufio.NewScanner(r.Body)
	reader.Scan()
	rawurl := reader.Text()

	// Ensure url given is a real url
	cleanUrl, err := CleanURL(rawurl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid url \"%s\"", reader.Text()), http.StatusBadRequest)
		return
	}

	// Generate short code and store in cache
	code := GetShortCode([]byte(cleanUrl))
	handler.cache.Add(code, cleanUrl)

	// Generate response url
	codeUrl := &url.URL{
		Scheme: "https",
		Host:   r.Host,
		Path:   "/" + code,
	}
	fmt.Fprintf(w, codeUrl.String())
}

func NewServer(bind string, maxEntries int) *http.Server {
	return &http.Server{
		Addr: bind,
		Handler: &Handler{
			cache: lru.New(maxEntries),
		},
	}
}
