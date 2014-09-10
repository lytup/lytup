package lytup

import (
	"encoding/json"
	"io/ioutil"

	"github.com/golang/glog"
)

var (
	dirs   = [...]string{"/etc", "/usr/local/etc", "."} // Directories where config file can be found
	Config config
)

type mongoDbConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type redisConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
}

type emailConfig struct {
	Host      string `json:"host"`
	Port      uint   `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromName  string `json:"fromName"`
	FromEmail string `json:"fromEmail"`
}

type messageConfig struct {
	EmailIsRegisteredError string `json:"emailIsRegisteredError"`
	VerifyEmailSuccess     string `json:"verifyEmailSuccess"`
	VerifyEmailFailed      string `json:"verifyEmailFailed"`
	EmailIsVerifiedError   string `json:"emailIsVerifiedError"`
	EmailNotFoundError     string `json:"emailNotFoundError"`
	ResetPasswordSuccess   string `json:"resetPasswordSuccess"`
	ResetPasswordFailed    string `json:"resetPasswordFailed"`
	UserNotFoundError      string `json:"userNotFoundError"`
	LoginFailed            string `json:"loginFailed"`
	ValidateTokenFailed    string `json:"validateTokenFailed"`
}

type config struct {
	Hostname            string        `json:"hostname"`
	Key                 string        `json:"key"`
	VerifyEmailExpiry   uint          `json:"verifyEmailExpiry"`
	PasswordResetExpiry uint          `json:"PasswordResetExpiry"`
	UploadDirectory     string        `json:"uploadDirectory"`
	MongoDb             mongoDbConfig `json:"mongoDb"`
	Redis               redisConfig   `json:"redis"`
	Email               emailConfig   `json:"email"`
	Message             messageConfig `json:"message"`
}

func init() {
	// Load config
	for _, d := range dirs {
		b, err := ioutil.ReadFile(d + "/lytup.json")
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &Config); err != nil {
			glog.Fatal(err)
		}
		return
	}
	glog.Fatalf("Config file not found in %v", dirs)
}
