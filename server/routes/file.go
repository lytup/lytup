package routes

import (
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// TODO: Move it to config file
const (
	UPLOAD_DIR = "/tmp"
)

func CreateFile(params martini.Params, ren render.Render, file models.File) {
	file.Create(params["id"])
	ren.JSON(http.StatusCreated, file)
}

func FindFileById(params martini.Params, ren render.Render) {
	_, file := models.FindFileById(params["id"])
	ren.JSON(http.StatusOK, file)
}

func UpdateFile(rw http.ResponseWriter, params martini.Params, file models.File) {
	models.UpdateFile(params["folId"], params["id"], &file)
	rw.WriteHeader(http.StatusOK)
}

func UploadFiles(req *http.Request, rw http.ResponseWriter, params martini.Params) {
	log.Println("Upload files")
	mr, err := req.MultipartReader()
	if err != nil {
		log.Fatal(err)
	}

	// Create folder
	folPath := path.Join(UPLOAD_DIR, params["id"])
	err = os.MkdirAll(folPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Get the file part
	part, err := mr.NextPart()
	if err != nil {
		log.Fatal(err)
	}
	defer part.Close()

	// Create file
	filePath := path.Join(folPath, part.FileName())
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	io.Copy(file, part)

	rw.WriteHeader(http.StatusNoContent)
}
