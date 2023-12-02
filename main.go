package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var dir = flag.String("dir", ".", "project dir, defaults to cwd")
var addr = flag.String("addr", ":8080", "http server address, defaults to :8080")

func main() {
	flag.Parse()
	var err error
	cmd := flag.Arg(0)
	switch cmd {
	case "serve":
		err = serve(*dir, *addr)
	case "render":
		err = renderAll(*dir)
	default:
		log.Printf("no command render or serve")
	}
	if err != nil {
		log.Printf("error running %s:\n%v", cmd, err)
	}
}

func serve(dir, addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		switch len(path) {
		case 1:
			// we expect 532 or 532_kids
			file := path[0]
			if serid.MatchString(file) {
				_, ext, _ := strings.Cut(file, ".")
				sub := "html"
				if ext == "txt" {
					sub = "text"
				}
				http.ServeFile(w, r, filepath.Join(dir, sub, file))
				return
			}
		case 0:
			http.ServeFile(w, r, filepath.Join(dir, "html", "index.html"))
			return
		}
		http.Error(w, "Seite nicht gefunden.", 404)
	})
	return http.ListenAndServe(addr, nil)
}
