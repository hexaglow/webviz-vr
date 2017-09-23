package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"os"
)

var bindAddress string = ":8000"
var staticDirectory string = "static"
var api MajesticApi

func envWithDefault(envVar string, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		log.Printf("%s not present in environment; using default value '%s'\n", envVar, defaultValue)
		return defaultValue
	} else {
		return value
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	majesticApiKey := os.Getenv("MAJESTIC_API_KEY")
	if majesticApiKey == "" {
		log.Fatalln("Cannot start: MAJESTIC_API_KEY missing from env.")
	}
	bindAddress := envWithDefault("BIND_ADDRESS", ":8000")
	staticDirectory := envWithDefault("STATIC_DIRECTORY", "static")

	api = MajesticApi{apiKey:majesticApiKey}

	log.Printf("Starting server on %s.\n", bindAddress)

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/test", TestHandler)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDirectory)))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(bindAddress, r))
}