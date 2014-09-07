package routes

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-martini/martini"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

func CreateFile(rw http.ResponseWriter, params martini.Params, ren render.Render, file models.File,
	usr *models.User) {
	if err := file.Create(params["folId"], usr); err != nil {
		glog.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		ren.JSON(http.StatusCreated, file)
	}
}

func FindFileById(rw http.ResponseWriter, params martini.Params, ren render.Render) {
	id, ok := params["id"]
	if !ok {
		id = params["fileId"]
	}

	if file, _, err := models.FindFileById(id); err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, file)
	}
}

func UpdateFile(rw http.ResponseWriter, params martini.Params, file models.File,
	usr *models.User) {
	file.Id = params["fileId"]
	if err := file.Update(params["folId"], usr); err != nil {
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

func DeleteFile(rw http.ResponseWriter, params martini.Params, usr *models.User) {
	if err := models.DeleteFile(params["folId"], params["fileId"], usr); err != nil {
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

func Upload(req *http.Request, rw http.ResponseWriter, params martini.Params,
	ren render.Render, usr *models.User) {
	var (
		folId string
		file  = &models.File{}
	)

	mr, err := req.MultipartReader()
	if err != nil {
		panic(err)
	}

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		defer part.Close()

		if name := part.FormName(); name != "file" {
			var buf bytes.Buffer
			io.Copy(&buf, part)
			val := buf.String()

			if name == "folId" {
				folId = val
				// Verify if folder belongs to the user
				_, err := models.FindFolderById(folId)
				if err != nil {
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
			} else if name == "fileId" {
				file.Id = val
				file.Uri = "/d/" + file.Id
			}
		} else if fileName := part.FileName(); fileName != "" {
			// Create folder
			folPath := path.Join(L.Config.UploadDirectory, folId)
			err = os.MkdirAll(folPath, 0755)
			if err != nil {
				panic(err)
			}

			// Create file
			filePath := path.Join(folPath, fileName)
			f, err := os.Create(filePath)
			defer f.Close()
			if err != nil {
				panic(err)
			}

			io.Copy(f, part)

			// Create thumbnail
			file.Thumbnail, err = utils.CreateThumbnail(filePath, part.Header.Get("Content-Type"))

			if err := file.Update(folId, usr); err != nil {
				glog.Error(err)
				if err == mgo.ErrNotFound {
					rw.WriteHeader(http.StatusNotFound)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
		}
	}
	ren.JSON(http.StatusOK, file)
}
