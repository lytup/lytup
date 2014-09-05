package routes

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
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
		ren.JSON(http.StatusCreated, usr.Render())
	}
}

func ConfirmUser(rw http.ResponseWriter, ren render.Render, params martini.Params) {
	// Get user id
	r := L.Redis()
	defer r.Close()
	id, err := redis.String(r.Do("GET", params["key"]))
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if usr, err := models.FindUserById(id); err != nil {
		glog.Error(err)
		if err == mgo.ErrNotFound {
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
			} else {
				ren.JSON(http.StatusOK, usr.Render())
			}
		}
	}
}

func ForgotPassword(rw http.ResponseWriter, usr *models.User) {

}

func FindUser(rw http.ResponseWriter, ren render.Render, usr *models.User) {
	if err := usr.Find(); err != nil {
		glog.Warning(err)
		if err == mgo.ErrNotFound {
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
		glog.Warning(err)
		if err == mgo.ErrNotFound {
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
		} else {
			glog.Warning(err)
		}
	}
	rw.WriteHeader(http.StatusUnauthorized)
}
