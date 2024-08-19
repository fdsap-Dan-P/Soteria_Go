// Package connect provides ...
package config

import (
	"fmt"
	"soteria_go/pkg/utils"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/encryptDecrypt"
	httpUtils "soteria_go/pkg/utils/go-utils/http"

	"log"
	"net/http"
)

func CreateConnection() {
	username := encryptDecrypt.DecodeBase64(utils.GetEnv("POSTGRES_USERNAME"))
	password := encryptDecrypt.DecodeBase64(utils.GetEnv("POSTGRES_PASSWORD"))
	host := encryptDecrypt.DecodeBase64(utils.GetEnv("POSTGRES_HOST"))
	dbName := encryptDecrypt.DecodeBase64(utils.GetEnv("DATABASE_NAME"))

	fmt.Println("username: ", username)
	fmt.Println("password: ", password)
	fmt.Println("host: ", host)
	fmt.Println("dbName: ", dbName)
	httpUtils.Client.New(&http.Client{})
	database.PostgreSQLConnect(
		username,
		password,
		host,
		dbName,
		utils.GetEnv("POSTGRES_PORT"),
		utils.GetEnv("POSTGRES_SSL_MODE"),
		utils.GetEnv("POSTGRES_TIMEZONE"),
	)
	err := database.DBConn.AutoMigrate()

	if err != nil {
		log.Fatal(err.Error())
	}

}
