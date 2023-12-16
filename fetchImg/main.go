package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Images struct {
	Path string `json:"path"`
	Date string `json:"date"`
}

type Receive struct {
	Dir string `json:"dir"`
}

var tgtDir string

func sorter(imgs []Images) []Images {
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Date < imgs[j].Date
	})
	return imgs
}

func fetchImg(tgtDirText string) ([]Images, error) {
	var imgLen int
	var imgPath []string
	var imgDate []string
	tgtDir := os.DirFS(tgtDirText)
	err := fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) { // ignore permission denied errors
			return err
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		fileType := strings.ToLower(path)
		if strings.Contains(fileType, "webp") || strings.Contains(fileType, "svg") ||
			strings.Contains(fileType, "jpeg") || strings.Contains(fileType, "jpg") ||
			strings.Contains(fileType, "gif") || strings.Contains(fileType, "png") ||
			strings.Contains(fileType, "tiff") || strings.Contains(fileType, "bmp") {
			imgPath = append(imgPath, path)
			imgDate = append(imgDate, fileInfo.ModTime().Format("20060102150405"))
			imgLen++
		}

		return nil
	})
	imgData := make([]Images, imgLen)
	if err != nil {
		return imgData, err
	}

	for i, v := range imgPath {
		imgData[i].Path = v
	}
	for i, v := range imgDate {
		imgData[i].Date = v
	}

	return imgData, nil
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
    tgtDir = receivedData.Dir

    // Fetch images after receiving the directory path
    imgs, err := fetchImg(tgtDir)
    if err != nil {
        http.Error(w, "Error fetching images", http.StatusInternalServerError)
        return
    }
    imgs = sorter(imgs)

    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    if err := enc.Encode(&imgs); err != nil {
        http.Error(w, "Error encoding images", http.StatusInternalServerError)
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
    fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"dir\": \"../\"}' http://localhost:8080/")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
