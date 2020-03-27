package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/damilarelana/goUrlShortener"
	"github.com/pkg/errors"
)

// defines the error message handler
func errMsgHandler(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

// define flags
var yamlFilename *string = flag.String("yaml", "pathsData.yaml", "a yaml file containing path and mapped URL, in a 'question, answer' format per record line")

// define reader to
//   * read in the file data
// 	 * return a slice of bytes

// fileOpener()
//   * gets the file from the system
//   * opens it and gets it ready to be used
//	 * returns a pointer to the now open/ready file
func fileOpener(f *string) *os.File {
	// Prepare the file content to be read by opening it first, using os.Open()
	openedFile, err := os.Open(*f)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to open file: %s\n", *f))
		panic(err)
	}
	return openedFile
}

// fileReader()
//  * takes the pointer to the opened file
//  * use ioutil to read the file and return data as []bytes
//  * returns the function call as []bytes
func fileReader(f *os.File) []byte {
	data, err := ioutil.ReadFile(*f)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to read file: %s\n", *f))
	}
	return data
}

//
// yaml parser is defined within:
//    * `goURlShortner` as parseYAML()
//

// urlShortenerHomepage handler
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

// defaultMux defines the router Mux that:
//   * initializes a new Mux
//   * maps routes to handlers
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortenerHomePage)
	return mux
}

// define main function that:
//   * uses defaultMux()
//   * uses mapHandler from `goURlShortner` package
//   * uses yamlHandler from `goURlShortner` package
func main() {
	// initialize all flags
	flag.Parse()

	// create an instance of defaultMux()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := goUrlShortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	yaml := `
- path: /urlshort
  url: https://github.com/damilarelana/goUrlShortener
- path: /urlshort-final
  url: https://github.com/damilarelana/goUrlShortener/tree/master/main
`

	yamlHandler, err := goUrlShortener.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to parse the YAML: %s\n", err))
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	log.Fatal(errors.Wrap(http.ListenAndServe(":8080", yamlHandler), "Failed to start WebServer"))
}
