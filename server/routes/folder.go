package routes

import (
	"archive/zip"
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func SaveFolder(fol models.Folder, ren render.Render) {
	fol.Save()
	ren.JSON(http.StatusCreated, fol)
}

func FindFolders(ren render.Render) {
	folders := models.FindFolders()
	ren.JSON(http.StatusOK, folders)
}

func FindFolderById(rw http.ResponseWriter, params martini.Params, ren render.Render) {
	fol := models.FindFolderById(params["id"])
	ren.JSON(http.StatusOK, fol)
}

func UpdateFolder(fol models.Folder, rw http.ResponseWriter, params martini.Params) {
	models.UpdateFolder(params["id"], &fol)
	rw.WriteHeader(http.StatusOK)
}

func DownloadFolder(rw http.ResponseWriter, params martini.Params) {
	log.Println("download folder")
	zw := zip.NewWriter(rw)
	fol := models.FindFolderById(params["id"])

	for _, file := range fol.Files {
		fw, err := zw.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Open(path.Join(UPLOAD_DIR, fol.Id, file.Name))
		io.Copy(fw, f)
	}

	err := zw.Close()
	if err != nil {
		log.Fatal(err)
	}

	rw.WriteHeader(http.StatusOK)
}
