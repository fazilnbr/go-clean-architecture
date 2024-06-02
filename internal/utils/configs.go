package utils

type AppConfig struct {
	Postgres       Postgres `json:"mysql"`
	TokenSecretKey string   `json:"token_secret_key"`
}

type Postgres struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}
