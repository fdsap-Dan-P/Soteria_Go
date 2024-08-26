// Package utils provides ...
package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load("project-env-files/.env")

	if err != nil {
		fmt.Println("Error loading .env file")
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
