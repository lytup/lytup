package routes

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

func CreateFolder(fol models.Folder, rw http.ResponseWriter, ren render.Render, usr *models.User) {
	fol.UserId = usr.Id
	if err := fol.Create(); err != nil {
		glog.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		ren.JSON(http.StatusCreated, fol)
	}
}

func FindFolders(rw http.ResponseWriter, ren render.Render, usr *models.User) {
	if folders, err := models.FindFolders(usr); err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, folders)
	}
}

func FindFolderById(rw http.ResponseWriter, params martini.Params, ren render.Render) {
	if fol, err := models.FindFolderById(params["id"]); err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, fol)
	}
}

func UpdateFolder(rw http.ResponseWriter, ren render.Render, params martini.Params, fol models.Folder,
	usr *models.User) {
	fol.Id = params["id"]
	if err := fol.Update(usr); err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, fol)
	}
}

func DeleteFolder(rw http.ResponseWriter, params martini.Params,
	usr *models.User) {
	err := models.DeleteFolder(params["id"], usr)
	if err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusOK)
	}
}

func Download(rw http.ResponseWriter, params martini.Params) {
	id := params["id"]

	db := L.NewDb("folders")
	defer db.Session.Close()

	// Is folder?
	if err := downloadFolder(id, rw); err == nil {
		return
	}

	// Is file?
	if err := downloadFile(id, rw, false); err != nil {
		glog.Error(err)
		rw.WriteHeader(http.StatusNotFound)
	}
}

func DownloadThumbnail(rw http.ResponseWriter, params martini.Params) {
	id := params["id"]
	downloadFile(id, rw, true)
}

func downloadFolder(id string, rw http.ResponseWriter) error {
	fol, err := models.FindFolderById(id)
	if err != nil {
		return err
	}

	zw := zip.NewWriter(rw)

	for _, file := range fol.Files {
		fw, err := zw.Create(file.Name)
		if err != nil {
			return err
		}

		f, err := os.Open(path.Join(L.Config.UploadDirectory, fol.Id, file.Name))
		defer f.Close()
		if err != nil {
			return err
		}

		rw.Header().Set("Content-Disposition", "attachment; filename="+fol.Id)
		io.Copy(fw, f)
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return nil
}

func downloadFile(id string, rw http.ResponseWriter, thumbnail bool) error {
	file, folId, err := models.FindFileById(id)
	if err != nil {
		return err
	}

	folPath := path.Join(L.Config.UploadDirectory, folId)
	f, e := os.Open(path.Join(folPath, file.Name))
	if thumbnail {
		if utils.IsVideo(file.Type) {
			file.Name += ".jpg" // Fetch as image
		}
		f, e = os.Open(path.Join(folPath, "t", file.Name))
	}
	defer f.Close()
	if e != nil {
		glog.Error(e)
	}

	fi, e := f.Stat()
	if e != nil {
		glog.Error(e)
	}

	rw.Header().Set("Content-Disposition", "attachment; filename='"+file.Name+"'")
	rw.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	io.Copy(rw, f)

	return nil
}
