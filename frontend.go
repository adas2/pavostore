package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var store *DataStore

var (
	flagPort = flag.String("port", "8080", "Port to listen on")
	flagFile = flag.String("filename", "input_data/example.data", "File path for the data store entries")
)

func main() {
	http.HandleFunc("/get", handleGet)

	flag.Parse()

	log.Fatal(http.ListenAndServe(":"+*flagPort, nil))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	store = getInstance(*flagFile)
	value, ok := store.Get(key)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Key '%s' not found\n", key)
		return
	}

	fmt.Fprintf(w, "%s = %s\n", key, value)
}
