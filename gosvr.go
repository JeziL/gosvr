package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/JeziL/gosvr/utils"
	"github.com/gobuffalo/packr/v2"

	gosvr "github.com/JeziL/gosvr/server"
)

const _Version = "1.0.6"

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
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Local IPs:")
	for _, ip := range utils.LocalIPs() {
		log.Printf("\t%s", ip)
	}
	log.Printf("Serving HTTP on 0.0.0.0 port %s (http://localhost:%s/) ...", *port, *port)
	log.Fatal(server.ListenAndServe())
}
