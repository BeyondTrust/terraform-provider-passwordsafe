package config

import (
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var PS_API_KEY string = os.Getenv("PS_API_KEY")
var PS_ACCOUNT_NAME string = os.Getenv("PS_ACCOUNT_NAME")
var PS_URL string = os.Getenv("PS_URL")
var CERTIFICATE_PATH string = os.Getenv("CERTIFICATE_PATH")
var CERTIFICATE_NAME string = os.Getenv("CERTIFICATE_NAME")
var CERTIFICATE_PASSWORD string = os.Getenv("CERTIFICATE_PASSWORD")
