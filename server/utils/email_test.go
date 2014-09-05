package utils

import (
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

// func TestWelcome(t *testing.T) {
// 	usr := models.User{
// 		Id:    bson.NewObjectId(),
// 		Name:  "Lab Test",
// 		Email: "test@labstack.com",
// 	}
// 	if err := Welcome(usr); err != nil {
// 		t.Error(err)
// 	}
// }

func TestEmailConfirmation(t *testing.T) {
	m := map[string]string{
		"hostname": "localhost:1431",
		"name":     "Vishal Rana",
		"email":    "vr@labstack.com",
		"key":      RandomString(32),
	}
	if err := EmailConfirmation(m); err != nil {
		t.Error(err)
	}
}
