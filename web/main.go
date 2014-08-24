package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/go-martini/martini"
	L "github.com/labstack/lytup/server/lytup"
	"gopkg.in/mgo.v2/bson"
)

var WEB_ROUTES map[string]bool = map[string]bool{
	"home": true,
}

func main() {
	m := martini.Classic()

	m.Get("/:id", home)   // Folder check
	m.Get("/i/:id", home) // File check

	log.Fatal(http.ListenAndServe("localhost:3001", m))
}

func home(req *http.Request, rw http.ResponseWriter, params martini.Params) {
	id := params["id"]

	if WEB_ROUTES[id] {
		sendIndex(rw)
		return
	}

	db := L.NewDb("folders")
	defer db.Session.Close()

	qry := bson.M{"id": id}
	if strings.HasPrefix(req.URL.Path, "/i/") {
		qry = bson.M{"files.id": id}
	}
	n, err := db.Collection.Find(qry).Count()
	if err != nil {
		log.Fatal(err)
	}
	if n != 0 {
		sendIndex(rw)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
}

func sendIndex(rw http.ResponseWriter) {
	file, err := os.Open("public/index.html")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Go HTTP server puts Content-Type automatically
	// rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, file)
}
