package router

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	. "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	"gopkg.in/mgo.v2/bson"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		ctx.Next()
	}
}

func Auth() HandlerFunc {
	return func(ctx *Context) {
		evt := "auth"
		parts := strings.Fields(ctx.req.Header.Get("Authorization"))
		if len(parts) == 2 {
			token := parts[1]
			t, err := jwt.Parse(token, func(t *jwt.Token) ([]byte, error) {
				return []byte(C.Key), nil
			})
			if err != nil {
				log.WithFields(log.Fields{
					"event": evt,
				}).Error(err)
				ctx.Render500()
			} else {
				id := t.Claims["sub"].(string)
				if !t.Valid {
					log.WithFields(log.Fields{
						"event": evt,
						"user":  id,
					}).Warn(err)
					ctx.Render(http.StatusUnauthorized, NewHttpError(http.StatusUnauthorized, M["error401"]))
				} else {
					usr := models.User{Id: bson.ObjectIdHex(id)}
					ctx.User = &usr
				}
			}
		}
		ctx.Next()
	}
}
