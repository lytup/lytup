package routes

import (
	"archive/zip"
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/db"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"io"
	"labix.org/v2/mgo/bson"
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
	fol.Id = params["id"]
	fol.Update()
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
		panic(err)
	}
	if n != 0 {
		downloadFolder(id, rw)
		return
	}

	// Check if file
	n, err = db.Collection.Find(bson.M{"files.id": id}).Count()
	if err != nil {
		panic(err)
	}
	if n != 0 {
		downloadFile(id, rw, false)
		return
	}

	rw.WriteHeader(http.StatusNotFound)
}

func DownloadThumbnail(rw http.ResponseWriter, params martini.Params) {
	id := params["id"]
	downloadFile(id, rw, true)
}

func downloadFolder(id string, rw http.ResponseWriter) {
	zw := zip.NewWriter(rw)
	fol := models.FindFolderById(id)

	for _, file := range fol.Files {
		fw, err := zw.Create(file.Name)
		if err != nil {
			panic(err)
		}

		f, err := os.Open(path.Join(UPLOAD_DIR, fol.Id, file.Name))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rw.Header().Set("Content-Disposition", "attachment; filename="+fol.Id)
		io.Copy(fw, f)
	}

	err := zw.Close()
	if err != nil {
		panic(err)
	}
}

func downloadFile(id string, rw http.ResponseWriter, thumbnail bool) {
	folId, file := models.FindFileById(id)
	folPath := path.Join(UPLOAD_DIR, folId)

	f, err := os.Open(path.Join(folPath, file.Name))
	if thumbnail {
		if utils.IsVideo(file.Type) {
			file.Name += ".jpg" // Fetch as image
		}
		f, err = os.Open(path.Join(folPath, "t", file.Name))
	}
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Disposition", "attachment; filename='"+file.Name+"'")
	rw.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	io.Copy(rw, f)
}
