package utils

type AppConfig struct {
	Postgres       Postgres `json:"mysql"`
	TokenSecretKey string   `json:"token_secret_key"`
}

type Postgres struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}
