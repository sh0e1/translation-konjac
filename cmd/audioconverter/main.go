package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Audio Converter")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
