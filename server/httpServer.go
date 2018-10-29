package gosvr

import (
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

	"github.com/JeziL/gosvr/utils"
)

// SimpleHTTPServer handles all HTTP requests.
type SimpleHTTPServer struct {
	Root    string
	Box     packr.Box
	Version string
}

type aFile struct {
	URL          string
	Filename     string
	Size         string
	IsDir        bool
	IsSymlink    bool
	IsFile       bool
	IsSourceCode bool
}

type aDir struct {
	Path    string
	Items   []aFile
	Version string
}

func (h SimpleHTTPServer) ServeHTTP(w loggingResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(&w, r)
	case http.MethodPost:
		h.post(&w, r)
	case http.MethodDelete:
		h.delete(&w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	log.Printf("%s - - \"%s %s %s\" %d - ", r.RemoteAddr, r.Method, r.URL.String(), r.Proto, w.StatusCode)
}

func (h SimpleHTTPServer) absPath(filePath string) string {
	return path.Join(h.Root, filePath)
}

func (h SimpleHTTPServer) getFiles(filePath string) []aFile {
	var items []aFile
	files, err := ioutil.ReadDir(h.absPath(filePath))
	utils.CheckError(err)
	for _, f := range files {
		fileURL := path.Join(filePath, f.Name())
		item := aFile{
			URL:          fileURL,
			Filename:     f.Name(),
			IsDir:        false,
			IsFile:       true,
			IsSymlink:    false,
			IsSourceCode: false,
			Size:         "",
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
			} else if b, lang := utils.IsSourceCode(f.Name()); b {
				item.IsSourceCode = true
				item.IsFile = false
				item.URL += fmt.Sprintf("?code=1&lang=%s&view=code", lang)
			}
			item.Size = utils.ByteToString(f.Size())
		}
		items = append(items, item)
	}
	return items
}

func (h SimpleHTTPServer) serveSourceCode(w *loggingResponseWriter, r *http.Request, filePath string, contentLength int64) {
	absPath := h.absPath(filePath)
	f, err := ioutil.ReadFile(absPath)
	utils.CheckError(err)
	data := struct {
		Path        string
		Lang        string
		FileContent string
		FileSize    string
		Version     string
	}{
		Path:        filePath,
		Lang:        r.URL.Query().Get("lang"),
		FileContent: string(f),
		FileSize:    utils.ByteToString(contentLength),
		Version:     h.Version,
	}
	t, err := template.New("codeView").Parse(h.Box.String("templates/codeView.html"))
	utils.CheckError(err)
	w.WriteHeader(http.StatusOK)
	err = t.Execute(w.Writer, data)
	utils.CheckError(err)
}

func (h SimpleHTTPServer) get(w *loggingResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	filePath, err := url.QueryUnescape(filePath)
	utils.CheckError(err)
	if strings.HasPrefix(filePath, "/gosvrstatic/") && r.URL.Query().Get("internal") == "1" {
		http.StripPrefix("/gosvrstatic/", http.FileServer(h.Box)).ServeHTTP(w, r)
		return
	}
	absPath := h.absPath(filePath)
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if utils.IsDir(absPath) {
		items := h.getFiles(filePath)
		data := aDir{
			Path:    filePath,
			Items:   items,
			Version: h.Version,
		}
		t, err := template.New("fileList").Parse(h.Box.String("templates/fileList.html"))
		utils.CheckError(err)
		w.WriteHeader(http.StatusOK)
		err = t.Execute(w.Writer, data)
		utils.CheckError(err)
	} else {
		fi, err := os.Stat(absPath)
		utils.CheckError(err)
		contentLength := fi.Size()
		if r.URL.Query().Get("code") == "1" && r.URL.Query().Get("view") == "code" {
			h.serveSourceCode(w, r, filePath, contentLength)
		} else {
			mimeType := utils.GuessType(fi.Name())
			const rfc2822 = "Mon, 02 Jan 15:04:05 -0700 2006"
			lastModified := fi.ModTime().Format(rfc2822)
			w.Header().Set("Content-type", mimeType)
			w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
			w.Header().Set("Last-Modified", lastModified)
			f, err := ioutil.ReadFile(absPath)
			utils.CheckError(err)
			w.Write(f)
		}
	}
}

func (h SimpleHTTPServer) post(w *loggingResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	utils.CheckError(err)
	if r.MultipartForm.File == nil {
		utils.CheckError(http.ErrMissingFile)
	}
	fhs := r.MultipartForm.File["files"]
	var fileNames []string
	for _, fh := range fhs {
		f, err := fh.Open()
		utils.CheckError(err)
		filename := fh.Filename
		fileNames = append(fileNames, filename)
		absPath := path.Join(h.absPath(r.URL.String()), filename)
		fw, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0666)
		utils.CheckError(err)
		io.Copy(fw, f)
		f.Close()
	}
	resultPage, err := template.New("uploaded").Parse(h.Box.String("templates/uploaded.html"))
	utils.CheckError(err)
	data := struct {
		FileNames []string
		Referer   string
		Version   string
	}{
		FileNames: fileNames,
		Referer:   r.Header.Get("Referer"),
		Version:   h.Version,
	}
	w.WriteHeader(http.StatusOK)
	err = resultPage.Execute(w.Writer, data)
	utils.CheckError(err)
}

func (h SimpleHTTPServer) delete(w *loggingResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	utils.CheckError(err)
	fileURL := r.FormValue("name")
	absPath := h.absPath(fileURL)
	err = os.RemoveAll(absPath)
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Fatal(err)
	}
}
