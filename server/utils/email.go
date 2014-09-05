package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/golang/glog"
	"github.com/jordan-wright/email"
	L "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/templates"
)

var (
	cfg  = L.Config.Email
	auth = smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.Host,
	)
)

func EmailConfirmation(data map[string]string) error {
	var b bytes.Buffer
	if err := templates.ExecuteTemplate(&b, "confirmation.html", data); err != nil {
		glog.Error(err)
		return err
	}

	e := &email.Email{
		To:      []string{data["email"]},
		From:    fmt.Sprintf("%s <%s>", cfg.FromName, cfg.FromEmail),
		Subject: "Welcome to Lytup",
		HTML:    b.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	if err := e.Send(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), auth); err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
