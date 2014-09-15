package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/golang/glog"
	"github.com/jordan-wright/email"
	. "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/templates"
)

var (
	auth = smtp.PlainAuth(
		"",
		C.Email.Username,
		C.Email.Password,
		C.Email.Host,
	)
)

func EmailVerifyEmail(data map[string]string) error {
	var b bytes.Buffer
	if err := templates.ExecuteTemplate(&b, "verify_email.html", data); err != nil {
		glog.Error(err)
		return err
	}

	e := &email.Email{
		To:      []string{data["email"]},
		From:    fmt.Sprintf("%s <%s>", C.Email.FromName, C.Email.FromEmail),
		Subject: "Welcome to Lytup",
		HTML:    b.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	if err := e.Send(fmt.Sprintf("%s:%d", C.Email.Host, C.Email.Port), auth); err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func EmailPasswordReset(data map[string]string) error {
	var b bytes.Buffer
	if err := templates.ExecuteTemplate(&b, "reset_pwd.html", data); err != nil {
		glog.Error(err)
		return err
	}

	e := &email.Email{
		To:      []string{data["email"]},
		From:    fmt.Sprintf("%s <%s>", C.Email.FromName, C.Email.FromEmail),
		Subject: "Lytup password reset",
		HTML:    b.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	if err := e.Send(fmt.Sprintf("%s:%d", C.Email.Host, C.Email.Port), auth); err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
