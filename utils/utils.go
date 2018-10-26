package utils

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path"
)

func IsDir(filePath string) bool {
	f, err := os.Stat(filePath)
	CheckError(err)
	if err != nil {
		return false
	}
	return f.Mode().IsDir()
}

func GuessType(filename string) string {
	ext := path.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if isCode, _ := IsSourceCode(filename); isCode {
		mimeType = "text/plain; charset=utf-8"
	}
	return mimeType
}

func ByteToString(b int64) string {
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

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsSourceCode(filename string) (bool, string) {
	var langExtMap = map[string]string{
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
	var langNameMap = map[string]string{
		"makefile": "makefile",
		"Makefile": "makefile",
	}
	ext := path.Ext(filename)
	name := path.Base(filename)[0 : len(filename)-len(ext)]
	val, exists := langExtMap[ext]
	if exists == false {
		val, exists = langNameMap[name]
	}
	return exists, val
}
