package routes

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	"github.com/labstack/lytup/server/models"
	U "github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo/bson"
	"net/http"
	"strings"
)

// var (
// 	PRIVATE_KEY []byte
// 	PUBLIC_KEY  []byte
// )
//
// func init() {
// 	usr, _ := user.Current()
// 	PRIVATE_KEY, _ = ioutil.ReadFile(path.Join(usr.HomeDir, ".ssh/id_rsa"))
// 	PUBLIC_KEY, _ = ioutil.ReadFile(path.Join(usr.HomeDir, ".ssh/id_rsa.pub"))
// 	log.Println(string(PRIVATE_KEY))
// }

func CreateUser(rw http.ResponseWriter, ren render.Render, usr models.User) {
	usr.HashedPassword = U.HashPassword([]byte(usr.Password))
	if err := usr.Create(); err != nil {
		glog.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		ren.JSON(http.StatusCreated, usr.Render())
	}

}

func FindUser(rw http.ResponseWriter, ren render.Render, usr *models.User) {
	if err := usr.Find(); err != nil {
		glog.Error(err)
		if err.Error() == "not found" {
			rw.WriteHeader(http.StatusNotFound)
		}
	} else {
		ren.JSON(http.StatusOK, usr.Render())
	}
}

func Login(rw http.ResponseWriter, ren render.Render, usr models.User) {
	if err := usr.Login(); err != nil {
		glog.Error(err)
		if err.Error() == "not found" {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, usr.Render())
	}
}

func ValidateToken(req *http.Request, rw http.ResponseWriter, ctx martini.Context) {
	parts := strings.Fields(req.Header.Get("Authorization"))
	if len(parts) == 2 {
		token := parts[1]
		t, err := jwt.Parse(token, func(t *jwt.Token) ([]byte, error) {
			return U.KEY, nil
		})
		if err == nil && t.Valid {
			id := t.Claims["usr-id"].(string)
			usr := models.User{Id: bson.ObjectIdHex(id)}
			ctx.Map(&usr)
			return
		}
	}
	rw.WriteHeader(http.StatusUnauthorized)
}
