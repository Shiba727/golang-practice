package main

import (
	"fmt"
	"net/http"
	"os"
	urlshort "urlshortner"

	_ "github.com/lib/pq"
)

func main() {
	listenAddr := os.Getenv("HTTP_PORT")
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	// yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	// if err != nil {
	// 	panic(err)
	// }

	db, err := urlshort.InitStore()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbHandler, err := urlshort.DBHandler(mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting the server on :%s", listenAddr)
	http.ListenAndServe(":"+listenAddr, dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", banner)
	mux.HandleFunc("/current", showCurrentPathinDB)
	return mux
}

func banner(w http.ResponseWriter, r *http.Request) {
	banner := `
	 _   _ ____  _           _                _
	| | | |  _ \| |      ___| |__   ___  _ __| |_ _ __   ___ _ __
	| | | | |_) | |     / __| '_ \ / _ \| '__| __| '_ \ / _ \ '__|
	| |_| |  _ <| |___  \__ \ | | | (_) | |  | |_| | | |  __/ |
	 \___/|_| \_\_____| |___/_| |_|\___/|_|   \__|_| |_|\___|_|
	`
	fmt.Fprintln(w, banner)
}

func showCurrentPathinDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	pathURLs, err := urlshort.FetchCurrentPathFromDB()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "------- Current path-url list in DB ------\n\n")
	for _, p := range pathURLs {
		fmt.Fprintf(w, "%d %s %s\n", p.ID, p.Path, p.URL)
	}
}
