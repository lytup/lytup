package routes

import (
	"bytes"
	"github.com/go-martini/martini"
	"github.com/labstack/lytup/server/models"
	"github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"io"
	"net/http"
	"os"
	"path"
)

// TODO: Move it to config file
const (
	UPLOAD_DIR = "/tmp"
)

func CreateFile(params martini.Params, ren render.Render, file models.File) {
	file.Create(params["folId"])
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

func UpdateFile(rw http.ResponseWriter, params martini.Params, file models.File) {
	file.Id = params["fileId"]
	file.Update(params["folId"])
	rw.WriteHeader(http.StatusOK)
}

func DeleteFile(rw http.ResponseWriter, params martini.Params) {
	models.DeleteFile(params["folId"], params["fileId"])
	rw.WriteHeader(http.StatusOK)
}

func Upload(req *http.Request, params martini.Params, ren render.Render) {
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
			if err != nil {
				panic(err)
			}
			defer f.Close()

			io.Copy(f, part)

			// Create thumbnail
			file.Thumbnail, err = utils.CreateThumbnail(filePath, part.Header.Get("Content-Type"))

			file.Update(folId)
		}
	}

	ren.JSON(http.StatusOK, file)
}
