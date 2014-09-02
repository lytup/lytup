package routes

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
)

// TODO: Move it to config file
const (
	UPLOAD_DIR = "/tmp"
)

func CreateFile(params martini.Params, ren render.Render, file models.File,
	usr *models.User) {
	file.Create(params["folId"], usr)
	ren.JSON(http.StatusCreated, file)
}

func FindFileById(params martini.Params, ren render.Render) {
	id, ok := params["id"]
	if !ok {
		id = params["fileId"]
	}
	_, file := models.FindFileById(id)
	ren.JSON(http.StatusOK, file)
}

func UpdateFile(rw http.ResponseWriter, params martini.Params, file models.File,
	usr *models.User) {
	file.Id = params["fileId"]
	file.Update(params["folId"], usr)
	rw.WriteHeader(http.StatusOK)
}

func DeleteFile(rw http.ResponseWriter, params martini.Params, usr *models.User) {
	models.DeleteFile(params["folId"], params["fileId"], usr)
	rw.WriteHeader(http.StatusOK)
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
			folPath := path.Join(UPLOAD_DIR, folId)
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

			file.Update(folId, usr)
		}
	}

	ren.JSON(http.StatusOK, file)
}
