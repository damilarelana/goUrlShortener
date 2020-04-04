package goUrlShortener

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v3"
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
	// parse the YAML file
	pathUrls, err := parseYAML(yamlBytes)
	if err != nil {
		return nil, err
	}

	// convert parsedYAML into a map
	pathsToUrls := buildPathsMap(pathUrls)

	// re-use the MapHandler
	// * now return the newly padded pathsToUrls
	// * while returning it in a format that makes it look like you were calling MapHandler in the first place
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return an http.HandlerFunc (which also implements http.Handler)
//  * parse the JSON file
//  * convert parsedJSON into a map
//  * then re-use the MapHandler

// JSONHandler parses the JSON file [in byte form]
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the JSON file
	pathUrls, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}

	// convert parsedJSON into a map
	pathsToUrls := buildPathsMap(pathUrls)

	// re-use the MapHandler
	// * now return the newly padded pathsToUrls
	// * while returning it in a format that makes it look like you were calling MapHandler in the first place
	return MapHandler(pathsToUrls, fallback), nil
}

// SQLHandler return an http.HandlerFunc (which also implements http.Handler)
//  * accept the incoming slice of struct data from dbQuery() in main.go
//  * convert parsed SQL data into a map using buildPathsMap()
//  * then re-use the MapHandler

// SQLHandler parses the sql file [in byte form]
func SQLHandler(pathUrls []PathURL, fallback http.Handler) (http.HandlerFunc, error) {
	// convert pathUrls into a map
	pathsToUrls := buildPathsMap(pathUrls)

	// re-use the MapHandler
	// * now return the newly padded pathsToUrls
	// * while returning it in a format that makes it look like you were calling MapHandler in the first place
	return MapHandler(pathsToUrls, fallback), nil
}

// PathURL declares the type structure we'll parse the YAML or JSON or SQL data into
type PathURL struct {
	Path string `format:"path"`
	URL  string `format:"url"`
}

// buildPathsMap converts parsedYAML into a map i.e.
// * make empty map
// * fill up the empty map one at a time
// * using the data already parsed into `pathUrls`
// * i.e. for each `Path` there is a corresponding `URL`
func buildPathsMap(pTUrl []PathURL) map[string]string {
	pTUrls := make(map[string]string)
	for _, pu := range pTUrl {
		pTUrls[pu.Path] = pu.URL
	}
	return pTUrls
}

// parseYAML uses the `yaml` package to parse the YAML bytes into the Type struct pathURL
//  * yaml.Unmarshal reads `all` the content into memory at once
func parseYAML(yB []byte) (pathUrls []PathURL, err error) {
	err = yaml.Unmarshal(yB, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

// parseJSON uses the `json` package to parse the JSON bytes into the Type struct pathURL
//  * json.Unmarshal reads `all` the content into memory at once
func parseJSON(jB []byte) (pathUrls []PathURL, err error) {
	err = json.Unmarshal(jB, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}
