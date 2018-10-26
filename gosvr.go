package main

import (
	"flag"
	"github.com/JeziL/gosvr/server"
	"github.com/gobuffalo/packr"
	"log"
	"net/http"
	"time"
)

const _Version = "0.9.4"

func main() {
	var dir = flag.String("d", ".", "Root directory to serve files from.")
	var port = flag.String("p", "8080", "Port number of the HTTP service.")
	flag.Parse()

	box := packr.NewBox("./static")
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
