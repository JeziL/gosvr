package main

import (
	"flag"
	"github.com/gobuffalo/packr/v2"
	"log"
	"net/http"
	"time"

	"github.com/JeziL/gosvr/server"
)

const _Version = "1.0.0"

func main() {
	var dir = flag.String("d", ".", "Root directory to serve files from.")
	var port = flag.String("p", "8080", "Port number of the HTTP service.")
	flag.Parse()

	box := packr.New("gosvr", "./static")
	server := &http.Server{
		Addr:           ":" + *port,
		Handler:        gosvr.HTTPHandlerWrapper(gosvr.SimpleHTTPServer{Root: *dir, Box: box, Version: _Version}),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...", *port, *port)
	log.Fatal(server.ListenAndServe())
}
