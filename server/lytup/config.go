package lytup

import (
	"encoding/json"
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

type RedisConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
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
	Hostname            string        `json:"hostname"`
	Key                 string        `json:"key"`
	ConfirmationExpiry  uint          `json:"confirmationExpiry"`
	PasswordResetExpiry uint          `json:"PasswordResetExpiry"`
	UploadDirectory     string        `json:"uploadDirectory"`
	MongoDb             MongoDbConfig `json:"mongoDb"`
	Redis               RedisConfig   `json:"redis"`
	Email               EmailConfig   `json:"email"`
}

var Config config

func init() {
	// Load config
	for _, d := range dirs {
		b, err := ioutil.ReadFile(d + "/lytup.json")
		if err != nil {
			continue
		}
		json.Unmarshal(b, &Config)
		return
	}
	glog.Fatal(fmt.Sprintf("Config file not found in %v", dirs))
}
