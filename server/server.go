package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/routes"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"io"
	"log"
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

	//*********
	// Folders
	//*********
	m.Post("/api/folders", binding.Bind(models.Folder{}), routes.CreateFolder)
	m.Get("/api/folders", routes.FindFolders)
	m.Get("/api/folders/:id", routes.FindFolderById)
	m.Patch("/api/folders/:id", binding.Bind(models.Folder{}),
		routes.UpdateFolder)
	m.Delete("/api/folders/:id", routes.DeleteFolder)

	//*******
	// Files
	//*******
	m.Post("/api/folders/:id/files", binding.Bind(models.File{}),
		routes.CreateFile)
	m.Get("/api/folders/:folId/files/:id", routes.FindFileById)
	m.Get("/api/files/:id", routes.FindFileById)
	m.Patch("/api/folders/:folId/files/:id",
		binding.Bind(models.File{}), routes.UpdateFile)
	m.Delete("/api/folders/:folId/files/:id", routes.DeleteFile)

	//*******************
	// Upload / Download
	//*******************
	m.Post("/u/:id", routes.UploadFiles)
	m.Get("/d/:id", routes.Download)
	// Download files
	// https://github.com/visionmedia/express/blob/9bf1247716c1f43e2c31c96fc965387abfeae531/lib/utils.js#L161
	// m.Get("/d/i/:id", routes.DownloadFile)

	log.Fatal(http.ListenAndServe("localhost:3000", m))
}
