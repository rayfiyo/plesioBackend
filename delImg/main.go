package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Receive struct {
	Path string `json:"path"`
}

type Result struct {
	Status string `json:"code"`
}

var tgtPath string

func delImg(tgtPath string) (Result, error) {
	var result Result
	if err := os.Remove(tgtPath); err != nil {
		result.Status = "failure"
		return result, err
	}
	result.Status = "success"
	return result, nil
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var receivedData Receive
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receivedData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	tgtPath = receivedData.Path

	result, err := delImg(tgtPath)
	if err != nil {
		http.Error(w, "Error deleting images", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&result); err != nil {
		http.Error(w, "Error encoding result", http.StatusInternalServerError)
		return
	}

	fmt.Println(buf.String())
	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handlePostRequest)

	fmt.Println("http://localhost:8080/")
	fmt.Println("touch testImg/unneeded.png")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"path\": \"testImg/unneeded.png\"}' http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
