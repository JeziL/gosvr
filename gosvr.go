package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type simpleHTTPServer struct {
	Root string
}

type aFile struct {
	URL      string
	Filename string
}

type aDir struct {
	Path  string
	Items []aFile
}

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (h simpleHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s  \"%s %s %s\"", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
	t, err := template.ParseFiles("fileList.html")
	checkError(err)
	switch r.Method {
	case http.MethodGet:
		h.get(w, r, t)
	case http.MethodPost:
		h.post(w, r, t)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Unsupported method.")
	}
}

func (h simpleHTTPServer) absPath(filePath string) string {
	return path.Join(h.Root, filePath)
}

func (h simpleHTTPServer) getFiles(filePath string) []aFile {
	var items []aFile
	files, err := ioutil.ReadDir(h.absPath(filePath))
	checkError(err)
	for _, f := range files {
		url := path.Join(filePath, f.Name())
		item := aFile{
			URL:      url,
			Filename: f.Name(),
		}
		if f.IsDir() {
			item.Filename += "/"
		}
		items = append(items, item)
	}
	return items
}

func (h simpleHTTPServer) get(w http.ResponseWriter, r *http.Request, t *template.Template) {
	filePath := r.URL.String()
	absPath := h.absPath(filePath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		w.WriteHeader(404)
		return
	}
	if isDir(absPath) {
		items := h.getFiles(filePath)
		data := aDir{
			Path:  filePath,
			Items: items,
		}
		err := t.Execute(w, data)
		checkError(err)
	} else {
		fi, err := os.Stat(absPath)
		checkError(err)
		mimeType := guessType(path.Ext(absPath))
		contentLength := fi.Size()
		const rfc2822 = "Mon, 02 Jan 15:04:05 -0700 2006"
		lastModified := fi.ModTime().Format(rfc2822)
		w.Header().Set("Content-type", mimeType)
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
		w.Header().Set("Last-Modified", lastModified)
		f, err := ioutil.ReadFile(absPath)
		checkError(err)
		w.Write(f)
	}
}

func (h simpleHTTPServer) post(w http.ResponseWriter, r *http.Request, t *template.Template) {
	file, header, err := r.FormFile("file")
	checkError(err)
	defer file.Close()

	absPath := path.Join(h.absPath(r.URL.String()), header.Filename)
	f, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
	checkError(err)
	io.Copy(f, file)
	resultPage, err := template.ParseFiles("uploaded.html")
	checkError(err)
	data := struct {
		Filename string
		Referer  string
	}{
		Filename: absPath,
		Referer:  r.Header.Get("Referer"),
	}
	err = resultPage.Execute(w, data)
	checkError(err)
}

func main() {
	var dir = flag.String("d", ".", "Root directory to serve files from.")
	var port = flag.String("p", "8080", "Port number of the HTTP service.")
	flag.Parse()

	server := &http.Server{
		Addr:           ":" + *port,
		Handler:        &simpleHTTPServer{Root: *dir},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...", *port, *port)
	log.Fatal(server.ListenAndServe())
}
