package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/routes"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"io"
	"net/http"
)

// TODO: Move it to config file
const (
	UPLOAD_DIR = "/tmp"
)

type Message struct {
	Data string `json:"data"`
}

func wsHandler(ws *websocket.Conn) {
	// log.Println("websocket connected")
	io.Copy(ws, ws)
	// websocket.JSON.Send(ws, Message{"Hello"})
}

func main() {
	m := martini.Classic()

	m.Use(render.Renderer())

	// TODO: Fix me!
	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	m.Get("/ws", websocket.Handler(wsHandler).ServeHTTP)

	m.Group("/api", func(r martini.Router) {
		//*******
		// Users
		//*******
		r.Post("/users", binding.Bind(models.User{}), routes.CreateUser)
		r.Post("/users/login", binding.Bind(models.User{}), routes.Login)

		r.Get("/folders/:id", routes.FindFolderById)
	})

	m.Group("/api", func(r martini.Router) {
		//*******
		// Users
		//*******
		r.Get("/users", routes.FindUser)

		//*********
		// Folders
		//*********
		r.Post("/folders", binding.Bind(models.Folder{}), routes.CreateFolder)
		r.Get("/folders", routes.FindFolders)
		r.Patch("/folders/:id", binding.Bind(models.Folder{}),
			routes.UpdateFolder)
		r.Delete("/folders/:id", routes.DeleteFolder)

		//*******
		// Files
		//*******
		r.Post("/folders/:folId/files", binding.Bind(models.File{}),
			routes.CreateFile)
		r.Get("/folders/:folId/files/:fileId", routes.FindFileById)
		r.Get("/files/:id", routes.FindFileById)
		r.Patch("/folders/:folId/files/:fileId",
			binding.Bind(models.File{}), routes.UpdateFile)
		r.Delete("/folders/:folId/files/:fileId", routes.DeleteFile)
	}, routes.ValidateToken)

	//*******************
	// Upload / Download
	//*******************
	m.Post("/u", routes.ValidateToken, routes.Upload)
	m.Get("/d/:id", routes.Download)
	m.Get("/d/:id/t", routes.DownloadThumbnail)

	glog.Fatal(http.ListenAndServe("localhost:3000", m))
}
