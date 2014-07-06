package routes

import (
	"archive/zip"
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/db"
	"github.com/labstack/lytup/server/models"
	"github.com/martini-contrib/render"
	"io"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func CreateFolder(fol models.Folder, ren render.Render, usr *models.User) {
	fol.UserId = usr.Id
	fol.Create()
	ren.JSON(http.StatusCreated, fol)
}

func FindFolders(ren render.Render, usr *models.User) {
	folders := models.FindFolders()
	ren.JSON(http.StatusOK, folders)
}

func FindFolderById(params martini.Params, ren render.Render) {
	fol := models.FindFolderById(params["id"])
	ren.JSON(http.StatusOK, fol)
}

func UpdateFolder(ren render.Render, params martini.Params, fol models.Folder) {
	models.UpdateFolder(params["id"], &fol)
	ren.JSON(http.StatusOK, fol)
}

func DeleteFolder(rw http.ResponseWriter, params martini.Params) {
	models.DeleteFolder(params["id"])
	rw.WriteHeader(http.StatusOK)
}

func Download(rw http.ResponseWriter, params martini.Params) {
	id := params["id"]

	db := db.NewDb("folders")
	defer db.Session.Close()

	// Check if folder
	n, err := db.Collection.Find(bson.M{"id": id}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if n != 0 {
		downloadFolder(id, rw)
		return
	}

	// Check if file
	n, err = db.Collection.Find(bson.M{"files.id": id}).Count()
	if err != nil {
		log.Fatal(err)
	}
	if n != 0 {
		downloadFile(id, rw)
		return
	}

	rw.WriteHeader(http.StatusNotFound)
}

func downloadFolder(id string, rw http.ResponseWriter) {
	log.Println("Download folder")
	zw := zip.NewWriter(rw)
	fol := models.FindFolderById(id)

	for _, file := range fol.Files {
		fw, err := zw.Create(file.Name)
		if err != nil {
			log.Panic(err)
		}

		f, err := os.Open(path.Join(UPLOAD_DIR, fol.Id, file.Name))
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		rw.Header().Set("Content-Disposition", "attachment; filename="+fol.Id)
		io.Copy(fw, f)
	}

	err := zw.Close()
	if err != nil {
		log.Panic(err)
	}
}

func downloadFile(id string, rw http.ResponseWriter) {
	folId, file := models.FindFileById(id)
	folPath := path.Join(UPLOAD_DIR, folId)

	f, err := os.Open(path.Join(folPath, file.Name))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	rw.Header().Set("Content-Disposition", "attachment; filename='"+file.Name+"'")
	rw.Header().Set("Content-Length", strconv.FormatUint(file.Size, 10))
	io.Copy(rw, f)
}
