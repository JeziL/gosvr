package main

import (
	"fmt"
	"log"
	"mime"
	"os"
)

func isDir(filePath string) bool {
	f, err := os.Stat(filePath)
	checkError(err)
	if err != nil {
		return false
	}
	return f.Mode().IsDir()
}

func guessType(ext string) string {
	mimeType := mime.TypeByExtension(ext)
	plain := "text/plain; charset=utf-8"
	if mimeType == "" {
		var types = map[string]string{
			".c":  plain,
			".py": plain,
			".go": plain,
		}
		if val, exists := types[ext]; exists {
			mimeType = val
		}
	}
	return mimeType
}

func byteToString(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
