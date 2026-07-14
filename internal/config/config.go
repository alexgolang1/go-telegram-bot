package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	ApiURL   string
	DBDSN    string
}

func LoadConfig() *Config {

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("[FATAL] Ошибка загрузки config.env: %v", err)
	}

	token := os.Getenv("BOT_TOKEN")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if token == "" || host == "" || user == "" {
		log.Fatal("[FATAL] В config.env пропущены обязательные параметры!")
	}

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port + " sslmode=" + sslMode

	return &Config{
		BotToken: token,
		DBDSN:    dsn,
	}
}
