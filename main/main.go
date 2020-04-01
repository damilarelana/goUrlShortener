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
var yamlFilename *string = flag.String("yaml", "", "a yaml file containing path and mapped URL, in a 'question, answer' format per record line")
var jsonFilename *string = flag.String("json", "", "a json file containing path and mapped URL, in a 'question, answer' format per record line")

// fileReader()
//  * takes the pointer to the opened file
//	* dereferences the pointer to get the value i.e. *f
//  * use ioutil to read the file and return data as []bytes
//  * returns the function call as []bytes
func fileReader(f *string) []byte {
	data, err := ioutil.ReadFile(*f)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to read file: %s\n", *f))
	}
	return data
}

// multiFlagTester()
//  * checks to see if the user is trying to use multiple file formats (yaml, json, sql) at the same time
//  * check to set number of flags `n` that have been set
//	* throws an error message, if `n` > 1
//  * returns a true boolean, if `n` > 1
func multiFlagTester() bool {
	numFlag := flag.NFlag()
	if numFlag > 1 {
		return true
	}
	return false
}

// yamlFlagHandler()
func yamlFlagHandler(yamlFilename *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	yamlHandler, err := goUrlShortener.YAMLHandler(fileReader(yamlFilename), mapHandler)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to parse the YAML: %s\n", err))
		panic(err)
	}
	return yamlHandler
}

// jsonFlagHandler()
func jsonFlagHandler(jsonFilename *string, mapHandler http.HandlerFunc) http.HandlerFunc {
	jsonHandler, err := goUrlShortener.JSONHandler(fileReader(jsonFilename), mapHandler)
	if err != nil {
		errMsgHandler(fmt.Sprintf("Failed to parse the JSON: %s\n", err))
		panic(err)
	}
	return jsonHandler
}

// formatHandler()
// * validates to avoid multiple flags being used simultaneously
// * checks for which flag is being used
// * leverages the appropriate flag handler
// * defaults to using yaml flag when no flag is chosen
func selectFlagHandler(mapHandler http.HandlerFunc) http.HandlerFunc {
	if multiFlagTester() { // check if multiple flags are being used i.e. multiFlagTester returns a true or false boolean
		errMsgHandler(fmt.Sprintf("Cannot use multiple flags at once. Please choose only yaml or json or sql \n"))
	}

	if flag.NFlag() != 0 && !reflect.DeepEqual(jsonFilename, "") { // if json filename is used THEN yaml flag would be empty string
		jsonHandler := jsonFlagHandler(jsonFilename, mapHandler)
		return jsonHandler
	}

	if flag.NFlag() != 0 && !reflect.DeepEqual(yamlFilename, "") { // if yaml filename is used THEN json flag would be empty string
		yamlHandler := yamlFlagHandler(yamlFilename, mapHandler)
		return yamlHandler
	}

	// defaults to yamlHandler when no flags are selected i.e. 'flag.NFlag() returns 0'
	*yamlFilename = "pathsData.yaml"
	fmt.Printf("No file flags was set (using either '-yaml' or '-json'). Defaulting to file: %s\n", *yamlFilename)
	yamlHandler := yamlFlagHandler(yamlFilename, mapHandler)
	return yamlHandler
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
//   * uses jsonHandler from `goURlShortner` package
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
	fmt.Println("Starting the server on :8080")
	log.Fatal(errors.Wrap(http.ListenAndServe(":8080", selectFlagHandler(mapHandler)), "Failed to start WebServer"))
}
