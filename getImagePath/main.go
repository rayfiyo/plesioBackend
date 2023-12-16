package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

func main() {
	tgtDirText := "testImg"
	tgtDir := os.DirFS(tgtDirText)

	fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !os.IsPermission(err) {
				log.Fatal(err)
			}
		}
		if strings.Contains(path, ".png") {
			fmt.Println(path)
		}
		return nil
	})
}
