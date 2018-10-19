# gosvr

*A Golang alternative to SimpleHTTPServer, but more beautiful and powerful.*

![](https://i.loli.net/2018/10/18/5bc8638ff0897.png)

## Install

- **From source (with Go installed, recommended)**

	```bash
	go get github.com/JeziL/gosvr
	```

- **From binary**

	See [Latest release](https://github.com/JeziL/gosvr/releases/latest). 
	
	***!!!NOTE:** None of these prebuilt binaries have been tested except for **darwin/amd64**, **linux/amd64** and **windows/amd64**.*


## Usage

```
gosvr [-d DIR] [-p PORT]

optional arguments:
  -d DIR     Root directory to serve files from. (default ".")
  -p PORT    Port number of the HTTP service. (default "8080")
```

## Powered by

- [packr](https://github.com/gobuffalo/packr)
- [Bootstrap](https://getbootstrap.com/) \[[LICENSE](static/js/LICENSE)\] (along with [JQuery](https://jquery.com/) and [Popper.js](https://popper.js.org/) used by it)
- [Material Design icons by Google](https://github.com/google/material-design-icons) \[[LICENSE](static/assets/iconfont/LICENSE)\]
