package main

import (
	"fmt"
	"log"
	"net/http"
)

type Receive struct {
	Path string `json:"path"`
	Dir  string `json:"dir"`
}

type Result struct {
	Status string `json:"code"`
}

type Dirs struct {
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Depth    int    `json:"depth"`
}

type Images struct {
	Path string `json:"path"`
	Date string `json:"date"`
}

func main() {
	http.HandleFunc("/", HandlePostRequest)

	fmt.Println("http://localhost:8080/")
	fmt.Println("curl -X POST http://localhost:8080/")
	fmt.Println("touch testImg/unneeded.png")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"dir\": \"../\"}' http://localhost:8080/")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"path\": \"testImg/unneeded.png\"}' http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
