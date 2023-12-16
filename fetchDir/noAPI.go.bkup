package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

func main() {
	tgtDirText := "../"
	tgtDir := os.DirFS(tgtDirText)

	fs.WalkDir(tgtDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !os.IsPermission(err) {
				log.Fatal(err)
			}
		}
		if d.IsDir() == true {
			fmt.Println(path)
		}
		return nil
	})
}
