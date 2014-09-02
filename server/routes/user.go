package routes

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	U "github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	if err := usr.Create(); err != nil {
		glog.Error(err)
		data := map[string]interface{}{"error": err}
		if mgo.IsDup(err) {
			data["error"] = "duplicate"
		}
		ren.JSON(http.StatusInternalServerError, data)
	} else {
		// Send confirmation email
		m := map[string]string{
			"hostname": L.Config.Hostname,
			"name":     usr.Name,
			"email":    usr.Email,
			"code":     usr.ConfirmationCode,
		}
		go U.EmailConfirmation(m)

		ren.JSON(http.StatusCreated, usr.Render())
	}
}

func ConfirmUser(rw http.ResponseWriter, ren render.Render, params martini.Params) {
	if usr, err := models.FindUserByConfirmationCode(params["code"]); err != nil {
		glog.Error(err)
		if err.Code == L.MongoDbNotFoundError {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		if usr.Confirmed {
			// Already confirmed
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			usr.Confirmed = true
			if err := usr.Save(); err != nil {
				glog.Error(err)
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func FindUserById(rw http.ResponseWriter, ren render.Render, params martini.Params) {
	if usr, err := models.FindUserById(params["id"]); err != nil {
		glog.Error(err)
		if err.Code == L.MongoDbNotFoundError {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, usr.Render())
	}
}

func Login(rw http.ResponseWriter, ren render.Render, usr models.User) {
	if err := usr.Login(); err != nil {
		glob.Error(err)
		if err.Code == L.LoginError {
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
			return []byte(L.Config.Key), nil
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
