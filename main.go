package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		defaultPort := "3000"
		log.Println("no env var set for port, defaulting to " + defaultPort)

		http.Handle("/", http.FileServer(http.Dir("./static")))
		http.ListenAndServe(":"+defaultPort, nil)
	} else {
		http.Handle("/", http.FileServer(http.Dir("./static")))
		log.Println("starting server on port " + port)
		http.ListenAndServe(":"+port, nil)
	}
}
