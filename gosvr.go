package main

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const _Version = "0.9.1"

type simpleHTTPServer struct {
	Root string
	Box  packr.Box
}

type aFile struct {
	URL       string
	Filename  string
	Size      string
	IsDir     bool
	IsSymlink bool
	IsFile    bool
}

type aDir struct {
	Path    string
	Items   []aFile
	Version string
}

func (h simpleHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s  \"%s %s %s\"", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
	t, err := template.New("fileList").Parse(h.Box.String("templates/fileList.html"))
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
		fileUrl := path.Join(filePath, f.Name())
		item := aFile{
			URL:       fileUrl,
			Filename:  f.Name(),
			IsDir:     false,
			IsFile:    true,
			IsSymlink: false,
			Size:      "",
		}
		if f.IsDir() {
			item.Filename += "/"
			item.IsDir = true
			item.IsFile = false
		} else {
			if f.Mode()&os.ModeSymlink != 0 {
				item.Filename += "@"
				item.IsSymlink = true
				item.IsFile = false
			}
			item.Size = byteToString(f.Size())
		}
		items = append(items, item)
	}
	return items
}

func (h simpleHTTPServer) get(w http.ResponseWriter, r *http.Request, t *template.Template) {
	filePath := r.URL.String()
	filePath, err := url.QueryUnescape(filePath)
	checkError(err)
	if strings.HasPrefix(filePath, "/gosvrstatic/") {
		http.StripPrefix("/gosvrstatic/", http.FileServer(h.Box)).ServeHTTP(w, r)
		return
	}
	absPath := h.absPath(filePath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		w.WriteHeader(404)
		return
	}
	if isDir(absPath) {
		items := h.getFiles(filePath)
		data := aDir{
			Path:    filePath,
			Items:   items,
			Version: _Version,
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
	err := r.ParseMultipartForm(32 << 20)
	checkError(err)
	if r.MultipartForm.File == nil {
		checkError(http.ErrMissingFile)
	}
	fhs := r.MultipartForm.File["files"]
	var fileNames []string
	for _, fh := range fhs {
		f, err := fh.Open()
		checkError(err)
		filename := fh.Filename
		fileNames = append(fileNames, filename)
		absPath := path.Join(h.absPath(r.URL.String()), filename)
		fw, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
		checkError(err)
		io.Copy(fw, f)
		f.Close()
	}
	resultPage, err := template.New("uploaded").Parse(h.Box.String("templates/uploaded.html"))
	checkError(err)
	data := struct {
		FileNames []string
		Referer   string
		Version   string
	}{
		FileNames: fileNames,
		Referer:   r.Header.Get("Referer"),
		Version:   _Version,
	}
	err = resultPage.Execute(w, data)
	checkError(err)
}

func main() {
	var dir = flag.String("d", ".", "Root directory to serve files from.")
	var port = flag.String("p", "8080", "Port number of the HTTP service.")
	flag.Parse()

	box := packr.NewBox("./static")
	server := &http.Server{
		Addr:           ":" + *port,
		Handler:        &simpleHTTPServer{Root: *dir, Box: box},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...", *port, *port)
	log.Fatal(server.ListenAndServe())
}
