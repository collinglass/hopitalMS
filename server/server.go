package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Starting Server")
	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("./stub/"))))
	http.Handle("/", http.FileServer(http.Dir("../app/")))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
