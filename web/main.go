package main

import (
	"io"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
)

var box *rice.Box

func main() {
	m := martini.Classic()
	box = rice.MustFindBox("public")

	// Route everything to index page
	m.Get("/**", func(rw http.ResponseWriter) {
		file, err := box.Open("index.html")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Go HTTP server puts Content-Type automatically
		// rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		io.Copy(rw, file)
	})
	glog.Fatal(http.ListenAndServe("localhost:3001", m))
}
