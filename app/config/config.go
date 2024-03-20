package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_USER             string
	DB_PASSWORD         string
	DB_NAME             string
	MYSQL_DATABASE      string
	MYSQL_ROOT_PASSWORD string
	TZ                  string
	MYSQL_HOST          string
)

func LoadEnv() {
	err := godotenv.Load("env/api.env")

	if err != nil {
		log.Printf("読み込み出来ませんでした: %v", err)
	}
	DB_USER = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME = os.Getenv("DB_NAME")
	MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE")
	MYSQL_ROOT_PASSWORD = os.Getenv("MYSQL_ROOT_PASSWORD")
	MYSQL_HOST = os.Getenv("MYSQL_HOST")
	TZ = os.Getenv("TZ")
}
