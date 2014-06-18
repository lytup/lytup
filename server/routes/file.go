package routes

import (
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
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

func UploadFiles(req *http.Request, rw http.ResponseWriter, params martini.Params) {
	mr, err := req.MultipartReader()
	if err != nil {
		log.Fatal(err)
	}

	// Create folder
	folPath := path.Join(UPLOAD_DIR, params["folderId"])
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

func DownloadFiles(rw http.ResponseWriter, params martini.Params) {
	log.Println("Download files")
	folId := params["folderId"]
	fileId := params["fileId"]
	folPath := path.Join(UPLOAD_DIR, folId)

	file := models.FindFileById(folId, fileId)
	f, err := os.Open(path.Join(folPath, file.Name))
	if err != nil {
		log.Fatal(err)
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	io.Copy(rw, f)
}
