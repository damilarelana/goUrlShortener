package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/damilarelana/goUrlShortener"
	"github.com/pkg/errors"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := goUrlShortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: github.com/damilarelana/goUrlShortener
	// - path: /urlshort-final
	//   url: github.com/damilarelana/goUrlShortener/tree/solution
	// `
	// 	yamlHandler, err := goUrlShortener.YAMLHandler([]byte(yaml), mapHandler)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	fmt.Println("Starting the server on :8080")
	log.Fatal(errors.Wrap(http.ListenAndServe(":8080", mapHandler), "Failed to start WebServer"))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortenerHomePage)
	return mux
}

func urlShortenerHomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		custom404PageHandler(w, r, http.StatusNotFound)
		return
	}
	dataHomePage := "Url Shortener: homepage"
	io.WriteString(w, dataHomePage)
}

// custom404PageHandler defines custom 404 page
func custom404PageHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set("Content-Type", "text/html") // set the content header type
	w.WriteHeader(status)                       // this automatically generates a 404 status code
	if reflect.DeepEqual(status, http.StatusNotFound) {
		data404Page := "This page does not exist ... 404!" // custom error message content
		io.WriteString(w, data404Page)
	}
}
