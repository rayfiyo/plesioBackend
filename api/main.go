package main

import (
	"api/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.HandlePostRequest)

	fmt.Println("sample:")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"dir\": \"../\"}' http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
