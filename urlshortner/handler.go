package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	URL  string `yaml:"url,omitempty"`
	Path string `yaml:"path,omitempty"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			// redirect to the dest
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYml, err := parseYAML(yml)
	if err != nil {
		log.Fatal("Found err when parsing yaml: ", err)
	}
	pathMap := buildMap(parsedYml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(data []byte) ([]pathUrl, error) {
	var pathUrls []pathUrl
	err := yaml.Unmarshal(data, &pathUrls)
	return pathUrls, err
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathMap := make(map[string]string, len(pathUrls))
	for _, p := range pathUrls {
		pathMap[p.Path] = p.URL
	}
	return pathMap
}
