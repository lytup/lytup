package email

import (
	"bytes"
	"errors"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	"html/template"
)

const (
	fromName  string = "Lytup"
	fromEmail string = "no-reply@lytup.com"
)

type To struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Message struct {
	To        []To
	Subject   string
	FromName  string
	FromEmail string
	Text      string
	Html      string
}

type Email interface {
	Send(Message) error
}

func NewEmail() (Email, error) {
	e := L.Cfg.Email
	if e.Provider == "mandrill" {
		return NewMandrill(e), nil
	}
	return nil, errors.New("Email provider not configured")
}

func Welcome(usr models.User) error {
	// Read template
	t, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		return err
	}

	var b bytes.Buffer
	m := map[string]string{
		"name": usr.Name,
		"id":   usr.Id.Hex(),
	}
	err = t.Execute(&b, m)
	if err != nil {
		return err
	}

	email, err := NewEmail()
	if err != nil {
		glog.Fatal(err)
	}
	msg := Message{
		To:        []To{{Email: usr.Email, Name: usr.Email}},
		Subject:   "Welcome to Lytup",
		FromName:  fromName,
		FromEmail: fromEmail,
		Html:      b.String(),
	}
	err = email.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
