package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/lytup/server/models"
	R "github.com/labstack/lytup/server/router"
)

// import . "github.com/labstack/lytup/server/lytup"

func CreateFolder(ctx *R.Context) {
	evt := "create-folder"
	fol := &models.Folder{UserId: ctx.User.Id}
	if err := ctx.Bind(&fol); err != nil {
		log.Printf("xxx xxx %s", fol)
		handle500(ctx, err, evt, "")
	} else if err := fol.Create(); err != nil {

	} else {
		ctx.Render(http.StatusCreated, fol)
	}
	// fol.UserId = ctx.User.Id
	// if err := fol.Create(); err != nil {
	// 	glog.Error(err)
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// } else {
	// 	ren.JSON(http.StatusCreated, fol)
	// }
}

// func FindFolders(rw http.ResponseWriter, ren render.Render, usr *models.User) {
// 	if folders, err := models.FindFolders(usr); err != nil {
// 		glog.Error(err)
// 		if err == mgo.ErrNotFound {
// 			rw.WriteHeader(http.StatusNotFound)
// 		} else {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}
// 	} else {
// 		ren.JSON(http.StatusOK, folders)
// 	}
// }
//
// func FindFolderById(rw http.ResponseWriter, params martini.Params, ren render.Render) {
// 	if fol, err := models.FindFolderById(params["id"]); err != nil {
// 		glog.Error(err)
// 		if err == mgo.ErrNotFound {
// 			rw.WriteHeader(http.StatusNotFound)
// 		} else {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}
// 	} else {
// 		ren.JSON(http.StatusOK, fol)
// 	}
// }
//
// func UpdateFolder(rw http.ResponseWriter, ren render.Render, params martini.Params, fol models.Folder,
// 	usr *models.User) {
// 	fol.Id = params["id"]
// 	if err := fol.Update(usr); err != nil {
// 		glog.Error(err)
// 		if err == mgo.ErrNotFound {
// 			rw.WriteHeader(http.StatusNotFound)
// 		} else {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}
// 	} else {
// 		ren.JSON(http.StatusOK, fol)
// 	}
// }
//
// func DeleteFolder(rw http.ResponseWriter, params martini.Params,
// 	usr *models.User) {
// 	err := models.DeleteFolder(params["id"], usr)
// 	if err != nil {
// 		glog.Error(err)
// 		if err == mgo.ErrNotFound {
// 			rw.WriteHeader(http.StatusNotFound)
// 		} else {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 		}
// 	} else {
// 		rw.WriteHeader(http.StatusOK)
// 	}
// }
//
// func Download(rw http.ResponseWriter, params martini.Params) {
// 	id := params["id"]
//
// 	db := NewDb("folders")
// 	defer db.Close()
//
// 	// Is folder?
// 	if err := downloadFolder(id, rw); err == nil {
// 		return
// 	}
//
// 	// Is file?
// 	if err := downloadFile(id, rw, false); err != nil {
// 		glog.Error(err)
// 		rw.WriteHeader(http.StatusNotFound)
// 	}
// }
//
// func DownloadThumbnail(rw http.ResponseWriter, params martini.Params) {
// 	id := params["id"]
// 	downloadFile(id, rw, true)
// }
//
// func downloadFolder(id string, rw http.ResponseWriter) error {
// 	fol, err := models.FindFolderById(id)
// 	if err != nil {
// 		return err
// 	}
//
// 	zw := zip.NewWriter(rw)
//
// 	for _, file := range fol.Files {
// 		fw, err := zw.Create(file.Name)
// 		if err != nil {
// 			return err
// 		}
//
// 		f, err := os.Open(path.Join(C.UploadDirectory, fol.Id, file.Name))
// 		defer f.Close()
// 		if err != nil {
// 			return err
// 		}
//
// 		rw.Header().Set("Content-Disposition", "attachment; filename="+fol.Id)
// 		io.Copy(fw, f)
// 	}
//
// 	if err := zw.Close(); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func downloadFile(id string, rw http.ResponseWriter, thumbnail bool) error {
// 	file, folId, err := models.FindFileById(id)
// 	if err != nil {
// 		return err
// 	}
//
// 	folPath := path.Join(C.UploadDirectory, folId)
// 	f, e := os.Open(path.Join(folPath, file.Name))
// 	if thumbnail {
// 		if utils.IsVideo(file.Type) {
// 			file.Name += ".jpg" // Fetch as image
// 		}
// 		f, e = os.Open(path.Join(folPath, "t", file.Name))
// 	}
// 	defer f.Close()
// 	if e != nil {
// 		glog.Error(e)
// 	}
//
// 	fi, e := f.Stat()
// 	if e != nil {
// 		glog.Error(e)
// 	}
//
// 	rw.Header().Set("Content-Disposition", "attachment; filename='"+file.Name+"'")
// 	rw.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
// 	io.Copy(rw, f)
//
// 	return nil
// }

// func handleError(ctx *R.Context, err error, evt, id string) {
// 	log.WithFields(log.Fields{
// 		"event": evt,
// 		"user":  id,
// 	}).Warn(err)
//
// 	if _, ok := err.(ErrUserField); ok {
// 		ctx.Render(http.StatusBadRequest, NewHttpError(http.StatusBadRequest, err.Error()))
// 	} else {
// 		switch err {
// 		case ErrEmailIsRegistered:
// 			ctx.Render(http.StatusConflict, NewHttpError(http.StatusConflict, err.Error()))
// 		case ErrVerifyEmail, ErrEmailIsVerified, ErrResetPassword:
// 			ctx.Render(http.StatusBadRequest, NewHttpError(http.StatusBadRequest, err.Error()))
// 		case ErrEmailNotFound, ErrUserNotFound, ErrLogin:
// 			ctx.Render(http.StatusNotFound, NewHttpError(http.StatusNotFound, err.Error()))
// 		default:
// 			handle500(ctx, err, evt, id)
// 		}
// 	}
// }
