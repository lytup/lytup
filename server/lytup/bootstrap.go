package lytup

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"

	"github.com/GeertJohan/go.rice"
)

var (
	C Config
	M map[string]string
)

func init() {
	box := rice.MustFindBox("../")

	// Load config
	b := box.MustBytes("config.json")
	if err := json.Unmarshal(b, &C); err != nil {
		log.Fatal(err)
	}

	// Load message
	b = box.MustBytes("lytup/message.json")
	if err := json.Unmarshal(b, &M); err != nil {
		log.Fatal(err)
	}
}
