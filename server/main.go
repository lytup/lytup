package main

import (
	"io"

	"code.google.com/p/go.net/websocket"
	"github.com/labstack/lytup/server/handlers"
	R "github.com/labstack/lytup/server/router"
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
	r := R.New()
	r.Use(R.Logger())

	a := r.G("/api")
	sa := r.G("/api", R.Auth()) // Secured API

	// Users
	a.Post("/users", handlers.CreateUser)
	a.Post("/users/login", handlers.Login)
	a.Get("/users/verify/:key", handlers.VerifyEmail)
	a.Post("/users/forgot", handlers.ForgotPassword)
	a.Get("/users/reset/:key", handlers.ResetPassword)
	sa.Get("/users", handlers.FindUser)

	// Folders
	sa.Post("/folders", handlers.CreateFolder)
	// sa.Get("/folders", handlers.FindFolders)
	// sa.Patch("/folders/:id", handlers.UpdateFolder)
	// sa.Delete("/folders/:id", handlers.DeleteFolder)

	r.Run()
	// m := martini.Classic()
	//
	// m.Use(render.Renderer())
	//
	// // TODO: Fix me!
	// m.Use(cors.Allow(&cors.Options{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"*"},
	// 	AllowHeaders:     []string{"Content-Type"},
	// 	AllowCredentials: true,
	// }))
	//
	// m.Get("/ws", websocket.Handler(wsHandler).ServeHTTP)
	//
	// m.Group("/api", func(r martini.Router) {
	// 	//*******
	// 	// Users
	// 	//*******
	// 	r.Post("/users", binding.Bind(models.User{}), routes.CreateUser)
	// 	r.Post("/users/login", binding.Bind(models.User{}), routes.Login)
	// 	m.Get("/users/verify/:key", routes.VerifyEmail)
	// 	m.Post("/users/forgot", binding.Bind(models.User{}), routes.ForgotPassword)
	// 	m.Get("/users/reset/:key", routes.ResetPassword)
	//
	// 	r.Get("/folders/:id", routes.FindFolderById)
	// })
	//
	// m.Group("/api", func(r martini.Router) {
	// 	//*******
	// 	// Users
	// 	//*******
	// 	r.Get("/users", routes.FindUser)
	// 	// r.Patch("/users", binding.Bind(models.User{}), routes.UpdateUser)
	//
	// 	//*********
	// 	// Folders
	// 	//*********
	// 	r.Post("/folders", binding.Bind(models.Folder{}), routes.CreateFolder)
	// 	r.Get("/folders", routes.FindFolders)
	// 	r.Patch("/folders/:id", binding.Bind(models.Folder{}),
	// 		routes.UpdateFolder)
	// 	r.Delete("/folders/:id", routes.DeleteFolder)
	//
	// 	//*******
	// 	// Files
	// 	//*******
	// 	r.Post("/folders/:folId/files", binding.Bind(models.File{}),
	// 		routes.CreateFile)
	// 	r.Get("/folders/:folId/files/:fileId", routes.FindFileById)
	// 	r.Get("/files/:id", routes.FindFileById)
	// 	r.Patch("/folders/:folId/files/:fileId",
	// 		binding.Bind(models.File{}), routes.UpdateFile)
	// 	r.Delete("/folders/:folId/files/:fileId", routes.DeleteFile)
	// }, routes.ValidateToken)
	//
	// //*******************
	// // Upload / Download
	// //*******************
	// m.Post("/u", routes.ValidateToken, routes.Upload)
	// m.Get("/d/:id", routes.Download)
	// m.Get("/d/:id/t", routes.DownloadThumbnail)
	//
	// glog.Fatal(http.ListenAndServe("localhost:3000", m))
}
