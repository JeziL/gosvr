package main

import (
	"flag"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/JeziL/gosvr/server"
)

const _Version = "1.0.2"

func main() {
	var dir = flag.String("d", ".", "Root directory to serve files from.")
	var port = flag.String("p", "8080", "Port number of the HTTP service.")
	var version = flag.Bool("v", false, "Version number of gosvr.")
	flag.Parse()

	if *version {
		fmt.Printf("gosvr version gosvr%s %s/%s", _Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	box := packr.New("gosvr", "./static")
	server := &http.Server{
		Addr:           ":" + *port,
		Handler:        gosvr.HTTPHandlerWrapper(gosvr.SimpleHTTPServer{Root: *dir, Box: box, Version: _Version}),
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...", *port, *port)
	log.Fatal(server.ListenAndServe())
}
