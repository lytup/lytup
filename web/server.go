package main

import (
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/db"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
)

var WEB_ROUTES map[string]bool = map[string]bool {
	"home": true,
}

func home(rw http.ResponseWriter) {
	file, err := os.Open("public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Go HTTP server puts Content-Type automatically
	// rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, file)
}

func main() {
	m := martini.Classic()

	m.Get("/:id", func(rw http.ResponseWriter, req *http.Request, params martini.Params) {
		id := params["id"]

		if (WEB_ROUTES[id]) {
			home(rw)
			return
		}

		db := db.NewDb("folders")
		defer db.Session.Close()
		qry := db.Collection.Find(bson.M{"id": params["id"]})
		n, err := qry.Count()
		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			home(rw)
		}
	})

	log.Fatal(http.ListenAndServe("localhost:3001", m))
}
