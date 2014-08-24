package email

import (
	"bytes"
	"encoding/json"
	"errors"
	L "github.com/labstack/lytup/server/lytup"
	"net/http"
)

const uri string = "https://mandrillapp.com/api/1.0/messages/send.json"

type to struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type message struct {
	To        []to   `json:"to"`
	Subject   string `json:"subject"`
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Text      string `json:"text"`
	Html      string `json:"html"`
}

type Mandrill struct {
	Key     string  `json:"key"`
	Message message `json:"message"`
}

func NewMandrill(cfg L.EmailConfig) *Mandrill {
	return &Mandrill{
		Key: cfg.Key,
	}
}

func (md *Mandrill) Send(msg Message) error {
	md.Message.To = make([]to, len(msg.To))
	for i, t := range msg.To {
		md.Message.To[i].Name = t.Name
		md.Message.To[i].Email = t.Email
	}
	md.Message.Subject = msg.Subject
	md.Message.FromName = msg.FromName
	md.Message.FromEmail = msg.FromEmail
	md.Message.Text = msg.Text
	md.Message.Html = msg.Html

	body, err := json.Marshal(md)
	if err != nil {
		return err
	}

	res, err := http.Post(uri, "applicaiton/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		// Response returned as object
		var o map[string]interface{}
		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&o); err != nil {
			return err
		}
		if o["status"] == "error" {
			return errors.New(o["message"].(string))
		}
	} else {
		// Response returned as array
		var a []map[string]interface{}
		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&a); err != nil {
			return err
		}
		for _, o := range a {
			if o["status"] == "rejected" {
				return errors.New(o["reject_reason"].(string))
			}
		}
	}

	return nil
}
