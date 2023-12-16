package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

type Dirs struct {
	Path     string `json:"path"`
	FullPath string `json:"full_path"`
	Depth    int    `json:"depth"`
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

func main() {
	// fetchDir
	var dirs []Dirs
	dirs, err := fetchDir("./")
	if err != nil {
		log.Fatal(err)
	}

	// setup
	handler1 := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		if err := enc.Encode(&dirs); err != nil {
			log.Fatal(err)
		}
		fmt.Println(buf.String())

		_, err := fmt.Fprint(w, buf.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	// POST
	fmt.Println("http://localhost:8080/")
	fmt.Println("curl -X POST http://localhost:8080/")
	http.HandleFunc("/", handler1)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
