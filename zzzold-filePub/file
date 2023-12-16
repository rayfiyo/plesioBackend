package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func visit(path string, info os.FileInfo, err error) error {
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Permission denied for %s\n", path)
			return nil
		}

		return err
	}
	lines := strings.Split(path, "\n")
	for _, line := range lines {
		if strings.Contains(line, ".png") {
			fmt.Println(line)
		}
	}
	return nil
}

func main() {
	fmt.Printf("対象のディレクトリを入力:")
	targetPath := ""
	fmt.Scanf("%s", &targetPath)

	err := filepath.Walk(targetPath, visit)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
