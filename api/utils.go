package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ... (Receive, Result, Dirs, Images structs)

// fetchDir retrieves directory information.
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

// fetchImg retrieves image information.
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
