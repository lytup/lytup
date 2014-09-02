package lytup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
)

// Directories where config file can be found
var dirs = [...]string{"/etc", "/usr/local/etc", "."}

type MongoDbConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EmailConfig struct {
	Host      string `json:"host"`
	Port      uint   `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromName  string `json:"fromName"`
	FromEmail string `json:"fromEmail"`
}

type config struct {
	Hostname string        `json:"hostname"`
	Key      string        `json:"key"`
	MongoDb  MongoDbConfig `json:"mongoDb"`
	Email    EmailConfig   `json:"email"`
}

var Config config

func init() {
	if err := loadConfig(); err != nil {
		glog.Fatal(err)
	}
}

func loadConfig() error {
	for _, d := range dirs {
		b, err := ioutil.ReadFile(d + "/lytup.json")
		if err != nil {
			continue
		}
		json.Unmarshal(b, &Config)
		return nil
	}
	return errors.New(fmt.Sprintf("Config file not found in %v", dirs))
}
