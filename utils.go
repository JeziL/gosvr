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
	if isCode, _ := isSourceCode(ext); isCode {
		mimeType = "text/plain; charset=utf-8"
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

func isSourceCode(ext string) (bool, string) {
	var langMap = map[string]string{
		".bf":      "brainfuck",
		".c":       "cpp",
		".cc":      "cpp",
		".coffee":  "coffeescript",
		".conf":    "apache",
		".cpp":     "cpp",
		".cs":      "csharp",
		".css":     "css",
		".dart":    "dart",
		".diff":    "diff",
		".erl":     "erlang",
		".gemspec": "ruby",
		".go":      "golang",
		".h":       "cpp",
		".hpp":     "cpp",
		".hs":      "haskell",
		".ini":     "ini",
		".java":    "java",
		".js":      "javascript",
		".json":    "json",
		".jsp":     "java",
		".kt":      "kotlin",
		".lua":     "lua",
		".m":       "matlab",
		".md":      "markdown",
		".mm":      "objectivec",
		".php":     "php",
		".pl":      "perl",
		".plist":   "xml",
		".podspec": "ruby",
		".py":      "python",
		".r":       "r",
		".rb":      "ruby",
		".rs":      "rust",
		".scala":   "scala",
		".sh":      "bash",
		".sql":     "sql",
		".swift":   "swift",
		".tex":     "tex",
		".ts":      "typescript",
		".v":       "verilog",
		".vhdl":    "vhdl",
		".xml":     "xml",
		".yml":     "yaml",
	}
	val, exists := langMap[ext]
	return exists, val
}
