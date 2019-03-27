# gosvr

[![Build Status](https://travis-ci.com/JeziL/gosvr.svg?branch=master)](https://travis-ci.com/JeziL/gosvr) [![Go Report Card](https://goreportcard.com/badge/github.com/JeziL/gosvr)](https://goreportcard.com/report/github.com/JeziL/gosvr)

*A Golang alternative to SimpleHTTPServer, but more beautiful and powerful.*

<img src="https://i.loli.net/2018/12/16/5c15f9385ee73.jpg" /> <img src="https://i.loli.net/2018/12/16/5c15f9385f9f7.jpg" />

## Features

- Good-looking.
- View, download, upload, rename and delete files.
- Display source code with syntax highlighting.
- Render markdown files like GitHub does.

## Install

- **From source (with Go installed, recommended)**

	```bash
	go get github.com/JeziL/gosvr
	```

- **From binary (untested)**

	See [Latest release](https://github.com/JeziL/gosvr/releases/latest). 
	

## Usage

```
gosvr [-d DIR] [-p PORT]

optional arguments:
  -d DIR     Root directory to serve files from. (default ".")
  -p PORT    Port number of the HTTP service. (default "8080")
  -v         Version number of gosvr
  -h         Help info
```

## Powered by

- [packr](https://github.com/gobuffalo/packr) \[[LICENSE](https://github.com/gobuffalo/packr/blob/master/LICENSE.txt)\]
- [Bootstrap](https://getbootstrap.com/) \[[LICENSE](https://github.com/twbs/bootstrap/blob/v4-dev/LICENSE)\]
- [JQuery](https://jquery.com/) \[[LICENSE](https://github.com/jquery/jquery/blob/master/LICENSE.txt)\]
- [Popper.js](https://popper.js.org/) \[[LICENSE](https://github.com/FezVrasta/popper.js/blob/master/LICENSE.md)\]
- [highlight.js](https://highlightjs.org/) \[[LICENSE](https://github.com/highlightjs/highlight.js/blob/master/LICENSE)\]
- [highlightjs-line-numbers.js](https://github.com/wcoder/highlightjs-line-numbers.js/) \[[LICENSE](https://github.com/wcoder/highlightjs-line-numbers.js/blob/master/LICENSE)\]
- [showdown](https://github.com/showdownjs/showdown) \[[LICENSE](https://github.com/showdownjs/showdown/blob/master/LICENSE)\]
- [github-markdown-css](https://github.com/sindresorhus/github-markdown-css) \[[LICENSE](https://github.com/sindresorhus/github-markdown-css/blob/gh-pages/license)\]
- [Material Design icons by Google](https://github.com/google/material-design-icons) \[[LICENSE](https://github.com/google/material-design-icons/blob/master/LICENSE)\]
