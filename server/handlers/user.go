package handlers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	. "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	R "github.com/labstack/lytup/server/router"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2/bson"
)

func CreateUser(ctx *R.Context) {
	evt := "create-user"
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		handle500(ctx, err, evt, "")
	} else if err := usr.Create(); err != nil {
		handleError(ctx, err, evt, "")
	} else {
		ctx.Render(http.StatusCreated, usr.ToRender())
	}
}

func VerifyEmail(ctx *R.Context) {
	evt := "verify-email"
	r := Redis()
	defer r.Close()
	id, err := redis.String(r.Do("GET", "verify:email:"+ctx.P("code")))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrVerifyEmail
		}
		handleError(ctx, err, evt, id)
		return
	}

	usr := &models.User{Id: bson.ObjectIdHex(id)}
	if err := usr.FindById(); err != nil {
		if err == ErrUserNotFound {
			err = ErrVerifyEmail
		}
		handleError(ctx, err, evt, id)
		return
	}

	if usr.EmailVerified {
		// Already verified
		err = ErrEmailIsVerified
		handleError(ctx, err, evt, id)
	} else {
		usr.EmailVerified = true
		if err := usr.Save(); err != nil {
			handleError(ctx, err, evt, id)
		} else {
			ctx.RenderOk(M["verifyEmailSuccess"])
		}
	}
}

func ForgotPassword(ctx *R.Context) {
	evt := "forgot-password"
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		handle500(ctx, err, evt, "")
	} else if err := usr.FindByEmail(); err != nil {
		handleError(ctx, err, evt, usr.Email)
	} else {
		// Create password reset key with expiry
		code := U.RandomString(32)
		key := "reset:pass:" + code
		val := usr.Id.Hex()
		if err := SetKeyWithExpiry(key, val, C.PasswordResetExpiry); err != nil {
			handle500(ctx, err, evt, "")
			return
		}

		// Send password reset email
		m := map[string]string{
			"name":  usr.FirstName,
			"email": usr.Email,
			"code":  code,
		}
		go U.EmailPasswordReset(m)
	}
}

func ResetPassword(ctx *R.Context) {
	evt := "reset-password"
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		handle500(ctx, err, evt, "")
	} else {
		r := Redis()
		defer r.Close()
		id, err := redis.String(r.Do("GET", "reset:pass:"+ctx.P("code")))
		if err != nil {
			if err == redis.ErrNil {
				err = ErrResetPassword
			}
			handleError(ctx, err, evt, id)
		} else {
			usr.Id = bson.ObjectIdHex(id)
			if err := usr.ResetPassword(); err != nil {
				handleError(ctx, err, evt, id)
			} else {
				usr.Login(false)
				ctx.Render(http.StatusOK, usr.ToRender())
			}
		}
	}
}

func Login(ctx *R.Context) {
	evt := "login"
	var usr models.User
	if err := ctx.Bind(&usr); err != nil {
		handle500(ctx, err, evt, "")
	} else if err := usr.Login(true); err != nil {
		handleError(ctx, err, evt, usr.Email)
	} else {
		ctx.Render(http.StatusOK, usr.ToRender())
	}
}

// FindUser finds authorized user.
func FindUser(ctx *R.Context) {
	evt := "find-user"
	usr := ctx.User
	if err := usr.FindById(); err != nil {
		handleError(ctx, err, evt, usr.Id.Hex())
	} else {
		ctx.Render(http.StatusOK, usr.ToRender())
	}
}

// handleError logs error and calls Context.Render
func handleError(ctx *R.Context, err error, evt, id string) {
	log.WithFields(log.Fields{
		"event": evt,
		"user":  id,
	}).Warn(err)

	if _, ok := err.(ErrUserField); ok {
		ctx.Render(http.StatusBadRequest, NewHttpError(http.StatusBadRequest, err.Error()))
	} else {
		switch err {
		case ErrEmailIsRegistered:
			ctx.Render(http.StatusConflict, NewHttpError(http.StatusConflict, err.Error()))
		case ErrVerifyEmail, ErrEmailIsVerified, ErrResetPassword:
			ctx.Render(http.StatusBadRequest, NewHttpError(http.StatusBadRequest, err.Error()))
		case ErrEmailNotFound, ErrUserNotFound, ErrLogin:
			ctx.Render(http.StatusNotFound, NewHttpError(http.StatusNotFound, err.Error()))
		default:
			handle500(ctx, err, evt, id)
		}
	}
}

// handle500 logs error and calls Context.Render500
func handle500(ctx *R.Context, err error, evt, id string) {
	log.WithFields(log.Fields{
		"event": evt,
		"user":  id,
	}).Error(err)
	ctx.Render500()
}
