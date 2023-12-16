package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type Dirs struct {
	Path string `json:"path"`
}

func fetchDir(tgtDirText string) ([]Dirs, error) {
	var dirPath []string
	var dirLen int
	tgtDir := os.DirFS(tgtDirText)
	err := fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) { // ignore permission denied errors
			return err
		}
		if d.IsDir() == true {
			dirPath = append(dirPath, path)
			dirLen++
		}
		return nil
	})
	dirData := make([]Dirs, dirLen)
	if err != nil {
		return dirData, err
	}

	for i, v := range dirPath {
		dirData[i].Path = v
	}

	return dirData, nil
}

func main() {
	// fetchDir
	dirs := make([]Dirs, 1)
	dirs, err := fetchDir("../../")
	if err != nil {
		return
	}

	// setup
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&dirs); err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf.String())

		_, err := fmt.Fprint(w, buf.String())
		if err != nil {
			return
		}
	}

	// GET
	fmt.Println("http://localhost:8080/images")
	http.HandleFunc("/images", handler1)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
