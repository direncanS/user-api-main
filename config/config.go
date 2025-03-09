package config

import (
	"os"
)

var (
	Host          = os.Getenv("host")
	Password      = os.Getenv("password")
	Port          = 5432
	User          = "postgres"
	Database      = "postgres"
	SecretKey     = []byte(os.Getenv("SECRET_KEY"))
	DB            = os.Getenv("DB")
	Addr          = os.Getenv("endpoint")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
)
