package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
)

// GLOBALS
var dbx *sqlx.DB

// Add go generate directives here, but using // not //// to start line
////go:generate go run $PWD/src/templating/main.go -name "human_names" -kind HumanNames -interval 300 -looker "LLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL"

var cache map[string]refresher
var refreshChan chan string

func setupCaches() {
	/*
	   Each individual entry registers itself with the cache and creates cache if cache is empty.
	*/
}

func setup() {
	setupDb()
	setupCaches()
	setRefresher(cache)
	warmCache()
}

func warmCache() {
	for _, entry := range cache {
		entry.refetchData()
	}
}

func main() {
	setup()

	for _, e := range cache {
		http.HandleFunc(e.LookerRoute(), e.Handler)
		http.HandleFunc(e.NamedRoute(), e.Handler)
	}

	http.HandleFunc("/looker.json", func(w http.ResponseWriter, req *http.Request) {
		var routes = map[string][]string{}
		routes["routes"] = []string{}
		for _, entry := range cache {
			routes["routes"] = append(routes["routes"], entry.LookerRoute())
		}
		js, err := json.Marshal(routes)
		if err != nil {
			log.Printf("Unable to marshal json")
		}
		io.WriteString(w, string(js))
	})

	http.HandleFunc("/named.json", func(w http.ResponseWriter, req *http.Request) {
		var routes = map[string][]string{}
		routes["routes"] = []string{}
		for _, entry := range cache {
			routes["routes"] = append(routes["routes"], entry.NamedRoute())
		}
		js, err := json.Marshal(routes)
		if err != nil {
			log.Printf("Unable to marshal json")
		}
		io.WriteString(w, string(js))
	})
	// DEBUG AND RESET CACHES
	http.HandleFunc("/debug", func(w http.ResponseWriter, req *http.Request) {
		go warmCache()
		io.WriteString(w, "ok")
	})
	env_port, err := os.LookupEnv("PORT")
	if err {
		env_port = "5000"
	}
	port := ":" + env_port
	log.Printf("Starting MViews Endpoint on %s\n", port)
	http.ListenAndServe(port, nil)
}
