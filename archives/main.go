package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func fetchDir(tgtDirText string) ([]Dirs, error) {
	var dirCount int
	var dirPath []string
	var dirDepth []int
	tgtDir := os.DirFS(tgtDirText)
	err := fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) { // permission denied は無視
			return err
		}
		if d.IsDir() == true {
			dirPath = append(dirPath, path)
			dirDepth = append(dirDepth, strings.Count(path, "/"))
			dirCount++
		}
		return nil
	})
	dirData := make([]Dirs, dirCount-1) // 先頭のカレントディレクトリは抜かす
	if err != nil {
		return dirData, err
	}

	if workingPath, errGetwd := os.Getwd(); errGetwd != nil {
		return dirData, errGetwd
	} else {
		workingPath = workingPath + "/"
		workingDepth := strings.Count(workingPath, "/")
		for i := 1; i < dirCount; i++ { // 先頭のカレントディレクトリは抜かす
			dirData[i-1].Path = dirPath[i]
			dirData[i-1].FullPath = workingPath + dirPath[i]
			dirData[i-1].Depth = workingDepth + dirDepth[i]
		}
	}

	return dirData, nil
}

func fetchImg(tgtDirText string) ([]Images, error) {
	var imgLen int
	var imgPath []string
	var imgDate []string
	tgtDir := os.DirFS(tgtDirText)
	err := fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) { // permission denied は無視
			return err
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		lowCharPath := strings.ToLower(path)
		if strings.Contains(lowCharPath, "webp") || strings.Contains(lowCharPath, "svg") ||
			strings.Contains(lowCharPath, "jpeg") || strings.Contains(lowCharPath, "jpg") ||
			strings.Contains(lowCharPath, "gif") || strings.Contains(lowCharPath, "png") ||
			strings.Contains(lowCharPath, "tiff") || strings.Contains(lowCharPath, "bmp") {
			imgPath = append(imgPath, filepath.ToSlash(path))
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

func main() {
	http.HandleFunc("/", handlePostRequest)

	fmt.Println("http://localhost:8080/")
	fmt.Println("curl -X POST http://localhost:8080/")
	fmt.Println("touch testImg/unneeded.png")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"dir\": \"../\"}' http://localhost:8080/")
	fmt.Println("curl -X POST -H \"Content-Type: application/json\" -d '{\"path\": \"testImg/unneeded.png\"}' http://localhost:8080/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
