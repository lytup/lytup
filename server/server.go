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
	// Create folder
	m.Post("/api/folders", binding.Bind(models.Folder{}), routes.SaveFolder)
	// Get folders
	m.Get("/api/folders", routes.FindFolders)
	// Get folder
	m.Get("/api/folders/:id", routes.FindFolderById)
	// Upload files
	m.Post("/u/:folderId", routes.UploadFiles)
	// Update folder
	m.Patch("/api/folders/:id", binding.Bind(models.Folder{}), routes.UpdateFolder)
	// Download folder
	m.Get("/d/:id", routes.DownloadFolder)
	// Download files
	// https://github.com/visionmedia/express/blob/9bf1247716c1f43e2c31c96fc965387abfeae531/lib/utils.js#L161
	m.Get("/d/:folderId/:fileId", routes.DownloadFiles)

	log.Fatal(http.ListenAndServe("localhost:3000", m))
}
