package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Audio Converter")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
