package email

import (
	"github.com/labstack/lytup/server/models"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

// func TestMandrillSend(t *testing.T) {
// 	email := NewMandrill()
// 	msg := Message{
// 		To:        []To{{Email: "test@labstack.com", Name: "Lab Test"}},
// 		Subject:   "Send test",
// 		FromName:  "Lab Test",
// 		FromEmail: "test@labstack.com",
// 		Html:      "<h1>Hello</h1>",
// 	}
// 	err := email.Send(msg)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestWelcome(t *testing.T) {
	usr := models.User{
		Id:    bson.NewObjectId(),
		Name:  "Lab Test",
		Email: "test@labstack.com",
	}
	err := Welcome(usr)
	if err != nil {
		t.Error(err)
	}
}
