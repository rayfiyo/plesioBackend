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

func fetchImg(tgtDirText string) ([]Images, error) {
	var imgLen int
	var imgPath []string
	var imgDate []string
	tgtDir := os.DirFS(tgtDirText)
	err := fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) { // ignore permission denied errors
			return err
		}
		if strings.Contains(path, ".png") {
			fileInfo, err := d.Info()
			if err != nil {
				return err
			}
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

func sorter(imgs []Images) []Images {
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Date < imgs[j].Date
	})
	return imgs
}

func main() {
	// fetchImage
	var imgs []Images
	imgs, err := fetchImg("../")
	if err != nil {
		log.Fatal(err)
	}

	imgs = sorter(imgs)

	// setup
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&imgs); err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf.String())

		_, err := fmt.Fprint(w, buf.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	// GET
	fmt.Println("http://localhost:8080/images")
	http.HandleFunc("/images", handler1)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
