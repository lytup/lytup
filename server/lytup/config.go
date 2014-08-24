package lytup

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
  "errors"
  "fmt"
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
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

type config struct {
	Key     string        `json:"key"`
	MongoDb MongoDbConfig `json:"mongoDb"`
	Email   EmailConfig   `json:"email"`
}

var Cfg config

func init() {
  err := loadConfig()
  if err != nil {
    glog.Fatal(err)
  }
}

func loadConfig() error {
  for _, d := range dirs {
    b, err := ioutil.ReadFile(d + "/lytup.json")
    if err != nil {
      continue
    }
    json.Unmarshal(b, &Cfg)
    return nil
  }
  return errors.New(fmt.Sprintf("Config file found in %v", dirs))
}
