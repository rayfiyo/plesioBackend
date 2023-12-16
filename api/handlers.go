package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

// ... (handlePostRequest, sorter, delImg functions)

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
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

	if receivedData.Path != "" {
		// delImg
		result, err := delImg(receivedData.Path)
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
	} else if receivedData.Dir != "" {
		// fetchImg
		imgs, err := fetchImg(receivedData.Dir)
		if err != nil {
			http.Error(w, "Error fetching images", http.StatusInternalServerError)
			return
		}
		imgs = sorter(imgs)

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&imgs); err != nil {
			http.Error(w, "Error encoding path of images", http.StatusInternalServerError)
			return
		}

		fmt.Println(buf.String())
		_, err = fmt.Fprint(w, buf.String())
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	}
}

func sorter(imgs []Images) []Images {
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Date < imgs[j].Date
	})
	return imgs
}

func delImg(tgtPath string) (Result, error) {
	var result Result
	if err := os.Remove(tgtPath); err != nil {
		result.Status = "failure"
		return result, err
	}
	result.Status = "success"
	return result, nil
}
