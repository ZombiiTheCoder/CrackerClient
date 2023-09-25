package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func root(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Clear-Size-Data", "*")
		fmt.Println(r.URL.Path)
		fs.ServeHTTP(w, r)
		
	}
}

func server(tdir, location, port string) {
	dir, _ := filepath.Abs(tdir)
	os.Chdir(dir)
	http.Handle("/", root(http.FileServer(http.Dir(dir))))

	log.Printf("Serving %s on HTTP port: %s\n", location, port)
	log.Fatal(http.ListenAndServe(location+":"+port, nil))
}

func serverEmbed(f embed.FS, location, port string) {
	d, _ := fs.Sub(f, "www")
	http.Handle("/", root(http.FileServer(http.FS(d))))

	log.Printf("Serving %s on HTTP port: %s\n", location, port)
	log.Fatal(http.ListenAndServe(location+":"+port, nil))
}