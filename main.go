package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
)

const BIND_ADDRESS string = ":8000"
const STATIC_DIRECTORY string = "static"
const STATIC_MOUNT_POINT string = "/"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	log.Printf("Starting server on %s.\n", BIND_ADDRESS)

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", RootHandler)

	r.PathPrefix(STATIC_MOUNT_POINT).Handler(
			http.StripPrefix(STATIC_MOUNT_POINT,
			http.FileServer(http.Dir(STATIC_DIRECTORY))))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(BIND_ADDRESS, r))
}