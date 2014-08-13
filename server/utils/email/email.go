package email

import (
	"bytes"
	"github.com/labstack/lytup/server/models"
	"html/template"
)

const (
	FROM_NAME  string = "Lytup"
	FROM_EMAIL string = "no-reply@lytup.com"
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

func NewEmail() Email {
	return NewMandrill()
}

func Welcome(usr models.User) error {
	// Read template
	t, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		return err
	}

	var b bytes.Buffer
	o := map[string]string{
		"name": usr.Name,
		"id":   usr.Id.Hex(),
	}
	err = t.Execute(&b, o)
	if err != nil {
		return err
	}

	email := NewEmail()
	msg := Message{
		To:        []To{{Email: usr.Email, Name: usr.Email}},
		Subject:   "Welcome to Lytup",
		FromName:  FROM_NAME,
		FromEmail: FROM_EMAIL,
		Html:      b.String(),
	}
	err = email.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
