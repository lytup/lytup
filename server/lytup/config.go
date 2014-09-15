package lytup

type (
	mongoDbConfig struct {
		Host     string `json:"host"`
		Port     uint   `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	redisConfig struct {
		Host     string `json:"host"`
		Port     uint   `json:"port"`
		Password string `json:"password"`
	}

	emailConfig struct {
		Host      string `json:"host"`
		Port      uint   `json:"port"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FromName  string `json:"fromName"`
		FromEmail string `json:"fromEmail"`
	}

	Config struct {
		Hostname            string        `json:"hostname"`
		Key                 string        `json:"key"`
		VerifyEmailExpiry   uint          `json:"verifyEmailExpiry"`
		PasswordResetExpiry uint          `json:"PasswordResetExpiry"`
		UploadDirectory     string        `json:"uploadDirectory"`
		MongoDb             mongoDbConfig `json:"mongoDb"`
		Redis               redisConfig   `json:"redis"`
		Email               emailConfig   `json:"email"`
	}
)
