package main

import (
	"fmt"
	"net/http"
	"os"
)

var version = "dev"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test-app version: %s\n", version)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("listening on :%s (version %s)\n", port, version)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
