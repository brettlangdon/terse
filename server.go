package terse

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/groupcache/lru"
)

type Handler struct {
	cache     *lru.Cache
	serverURL *url.URL
}

func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler.HandleGet(w, r)
	case "POST":
		handler.HandlePost(w, r)
	default:
		log.Printf("%s 405", r.Method)
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
		log.Printf("GET \"%s\" 301 \"%s\"", code, value.(string))
		http.Redirect(w, r, value.(string), http.StatusMovedPermanently)
		return
	}
	log.Printf("GET \"%s\" 404", code)
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
		log.Printf("POST \"%s\" 400", rawurl)
		http.Error(w, fmt.Sprintf("Invalid url \"%s\"", rawurl), http.StatusBadRequest)
		return
	}

	// Generate short code and store in cache
	code := GetShortCode([]byte(cleanUrl))

	// If the short code exists, and the urls are different, complain about the conflict
	value, ok := handler.cache.Get(code)
	if ok {
		if value.(string) != cleanUrl {
			log.Printf("POST \"%s\" 409", rawurl)
			msg := fmt.Sprintf("Short code conflict \"%s\" already registered as \"%s\"", code, value.(string))
			http.Error(w, msg, http.StatusConflict)
			return
		}
	} else {
		handler.cache.Add(code, cleanUrl)
	}

	// Generate response url
	codeUrl := handler.serverURL
	codeUrl.Path = "/" + code
	log.Printf("POST \"%s\" 201 created \"%s\"", rawurl, code)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, codeUrl.String())
}

func NewServer(bind string, maxEntries int, serverURL string) (*http.Server, error) {
	parsedURL, err := url.ParseRequestURI(serverURL)
	if err != nil {
		return nil, err
	}
	return &http.Server{
		Addr: bind,
		Handler: &Handler{
			cache:     lru.New(maxEntries),
			serverURL: parsedURL,
		},
	}, nil
}
