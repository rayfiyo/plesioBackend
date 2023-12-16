package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var imagePaths []string

type PageVariables struct {
	ImagePaths []string
}

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		if os.IsPermission(err) {
			return nil
		}
		return err
	}

	if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".png") {
		imagePaths = append(imagePaths, path)
	}

	return nil
}

func export() {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Image Gallery</title>
	</head>
	<body>
		<h1>Image Gallery</h1>
		{{range .ImagePaths}}
			<img src="{{.}}" alt="Image">
		{{end}}
	</body>
	</html>
	`

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pageVariables := PageVariables{
		ImagePaths: imagePaths,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = tmpl.Execute(w, pageVariables)
		if err != nil {
			fmt.Println("Error:", err)
		}
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	if dirs, err := os.ReadDir("./"); err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		for _, dir := range dirs {
			fmt.Println(dir)
		}
	}

	dirPath := "./"
	// fmt.Scanf(dirPath)

	if err := filepath.Walk(dirPath, visit); err != nil {
		fmt.Println("Error:", err)
		return
	}

	export()
}
