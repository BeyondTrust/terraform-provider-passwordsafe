package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var PS_API_KEY string = os.Getenv("PS_API_KEY")
var PS_ACCOUNT_NAME string = os.Getenv("PS_ACCOUNT_NAME")
var PS_URL string = os.Getenv("PS_URL")
