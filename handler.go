package goUrlShortener

import (
	"net/http"

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

// pathURL declares the type structure we'll parse the YAML into
type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

// buildPathsMap converts parsedYAML into a map i.e.
// * make empty map
// * fill up the empty map one at a time
// * using the data already parsed into `pathUrls`
// * i.e. for each `Path` there is a corresponding `URL`
func buildPathsMap(pTUrl []pathURL) map[string]string {
	pTUrls := make(map[string]string)
	for _, pu := range pTUrl {
		pTUrls[pu.Path] = pu.URL
	}
	return pTUrls
}

// parseYAML uses the `yaml` package to parse the YAML data
func parseYAML(yB []byte) (pathUrls []pathURL, err error) {
	err = yaml.Unmarshal(yB, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}
