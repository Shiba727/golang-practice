package main

import (
	"database/sql"
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
	yaml := `
    - path: /urlshort
      url: https://github.com/gophercises/urlshort
    - path: /urlshort-final
      url: https://github.com/gophercises/urlshort/tree/solution
    `
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	db, err := initStore()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("Starting the server on :%s", listenAddr)
	http.ListenAndServe(":"+listenAddr, yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", banner)
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

func initStore() (*sql.DB, error) {
	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	var (
		db  *sql.DB
		err error
	)

	db, err = sql.Open("postgres", pgConnString)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	// if _, err := db.Exec(
	// 	"CREATE TABLE IF NOT EXISTS message (value STRING PRIMARY KEY)"); err != nil {
	// 	return nil, err
	// }
	fmt.Println("Connected to DB.")
	return db, nil
}
