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
	U "github.com/labstack/lytup/server/utils"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2/bson"
)

var msg = L.Config.Message

func CreateUser(ren render.Render, usr models.User) {
	if err := usr.Create(); err != nil {
		if err == L.ErrEmailIsRegistered {
			handleError(ren, msg.EmailIsRegisteredError, http.StatusConflict)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusCreated, usr.Render())
	}
}

func VerifyEmail(params martini.Params, ren render.Render) {
	// Get user id
	r := L.Redis()
	defer r.Close()
	id, err := redis.String(r.Do("GET", "verify:email:"+params["key"]))
	if err != nil {
		if err == redis.ErrNil {
			handleError(ren, msg.VerifyEmailFailed, http.StatusNotFound)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if usr, err := models.FindUserById(id); err != nil {
		handleError(ren, err.Error(), http.StatusInternalServerError)
	} else {
		if usr.EmailVerified {
			// Already verified
			handleError(ren, msg.EmailIsVerifiedError, http.StatusBadRequest)
		} else {
			usr.EmailVerified = true
			if err := usr.Save(); err != nil {
				handleError(ren, err.Error(), http.StatusInternalServerError)
			} else {
				ren.JSON(http.StatusOK, map[string]string{
					"message": msg.VerifyEmailSuccess,
				})
			}
		}
	}
}

func ForgotPassword(ren render.Render, usr models.User) {
	if usr, err := models.FindUserByEmail(usr.Email); err != nil {
		if err == L.ErrEmailNotFound {
			handleError(ren, err.Error(), http.StatusNotFound)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Create password reset key with expiry
		r := L.Redis()
		defer r.Close()
		key := U.RandomString(32)
		k := "resetpwd:" + key
		r.Send("MULTI")
		r.Send("SET", k, usr.Id.Hex())
		r.Send("EXPIRE", k, L.Config.PasswordResetExpiry)
		if _, err := r.Do("EXEC"); err != nil {
			handleError(ren, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send password reset email
		m := map[string]string{
			"name":  usr.FirstName,
			"email": usr.Email,
			"key":   key,
		}
		go U.EmailPasswordReset(m)
	}
}

func ResetPassword(params martini.Params, ren render.Render) {
	// Get user id
	r := L.Redis()
	defer r.Close()
	id, err := redis.String(r.Do("GET", "resetpwd:"+params["key"]))
	if err != nil {
		if err == redis.ErrNil {
			handleError(ren, msg.ResetPasswordFailed, http.StatusNotFound)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if usr, err := models.FindUserById(id); err != nil {
		handleError(ren, err.Error(), http.StatusInternalServerError)
	} else {
		usr.Update()
		usr.Login(false)
		ren.JSON(http.StatusOK, usr.Render())
	}
}

func FindUser(ren render.Render, usr *models.User) {
	if err := usr.Find(); err != nil {
		if err == L.ErrUserNotFound {
			handleError(ren, err.Error(), http.StatusNotFound)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, usr.Render())
	}
}

func Login(ren render.Render, usr models.User) {
	if err := usr.Login(true); err != nil {
		if err == L.ErrLogin {
			handleError(ren, err.Error(), http.StatusNotFound)
		} else {
			handleError(ren, err.Error(), http.StatusInternalServerError)
		}
	} else {
		ren.JSON(http.StatusOK, usr.Render())
	}
}

// func UpdateUser(rw http.ResponseWriter, usr models.User) {
// 	if err := usr.Update(); err != nil {
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

func ValidateToken(req *http.Request, ren render.Render, ctx martini.Context) {
	parts := strings.Fields(req.Header.Get("Authorization"))
	if len(parts) == 2 {
		token := parts[1]
		t, err := jwt.Parse(token, func(t *jwt.Token) ([]byte, error) {
			return []byte(L.Config.Key), nil
		})
		if err != nil || !t.Valid {
			handleError(ren, msg.ValidateTokenFailed, http.StatusUnauthorized)
		} else {
			id := t.Claims["sub"].(string)
			usr := models.User{Id: bson.ObjectIdHex(id)}
			ctx.Map(&usr)
		}
	}
}

func handleError(ren render.Render, msg string, code int) {
	data := map[string]string{
		"message": msg,
	}
	if code == http.StatusInternalServerError {
		glog.Error(msg)
		data["message"] = "Looks like something went wrong!"
	} else {
		glog.Warning(msg)
	}
	ren.JSON(code, data)
}
