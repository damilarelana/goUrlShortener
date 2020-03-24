package goUrlShortener

import (
	"fmt"
	"net/http"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also implements http.Handler)
// * extract path in the request
// * find keys [in the map] that match the extracted path
// * redirect to the map value [for that key], if the key exists in the map
// * otherwise call the fallback http.Handler
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dest, ok := pathsToUrls[path]
		if ok { // `ok` would be true if `path` exists in pathsToUrls
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return an http.HandlerFunc (which also implements http.Handler)
//  * parse the YAML file
//  * convert parsedYAML into a map
//  * then re-use the MapHandler

// YAMLHandler parses the YAML file [in byte form]
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(yamlBytes, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert parsedYAML into a map
	// * make empty map
	// * fill up the empty map one at a time
	// * using the data already parsed into `pathUrls`
	// * i.e. for each `Path` there is a corresponding `URL`
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL

	}

	// re-use the MapHandler
	// * now return the newly padded pathsToUrls
	// * while returning it in a format that makes it look like you were calling MapHandler in the first place
	return MapHandler(pathsToUrls, fallback), nil
}

// middleware
type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// defines the error message handler
func errMsgHandler(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
