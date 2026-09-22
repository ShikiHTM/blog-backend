package auth

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	UserCredentialsFilePath string
	JWTSecretKey            string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, errors.New("Error loading .env file")
	}

	path, exist := os.LookupEnv("USER_CREDENTIALS_FILE_PATH")
	if !exist || path == "" {
		return Config{}, errors.New("USER_CREDENTIALS_FILE_PATH is required as a environment variable.")
	}

	secretKey, exist := os.LookupEnv("JWT_SECRET_KEY")
	if !exist || secretKey == "" {
		return Config{}, errors.New("JWT_SECRET_KEY is required as a environment variable.")
	}

	return Config{
		UserCredentialsFilePath: path,
		JWTSecretKey:            secretKey,
	}, nil
}
