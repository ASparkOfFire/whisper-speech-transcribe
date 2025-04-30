package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type config struct {
	ModelPath string `validate:"required"`
}

var (
	AppConfig config
)

func loadConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
		return err
	}

	AppConfig.ModelPath = os.Getenv("MODEL_PATH")

	return nil
}

func init() {
	if err := loadConfig(); err != nil {
		log.Println("Error setting application config: ", err)
		os.Exit(1)
	}
}
