package main

import (
	"log"
	"net/http"

	"psankar-goths-demo/handlers"
)

func main() {
	log.Println("Launched goths on :8080")
	http.HandleFunc("/", handlers.Login)

	http.ListenAndServe(":8080", nil)
}
