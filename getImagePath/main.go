package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PageVariables struct {
	ImagesPaths []string
}

func visit(variables *PageVariables) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), "png") {
			variables.ImagesPaths = append(variables.ImagesPaths, path)
		}
		return nil
	}
}

func main() {
	tgtDir := "testImg"

	var variables PageVariables
	if err := filepath.Walk(tgtDir, visit(&variables)); err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, v := range variables.ImagesPaths {
		fmt.Println(v)
	}
}
