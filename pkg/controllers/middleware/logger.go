package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func ActivityLogger(user string, endpointName string, retCode string, method string, endpoint string, requestBody []byte, responseBody []byte, message string, messagErr string, responsErr error) {
	// Set log directory
	logDir := "Logs"

	// Create log directory if it doesn't exist
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	// Create log directory for the endpoint
	endpointDir := filepath.Join(logDir, endpointName)
	if _, err := os.Stat(endpointDir); os.IsNotExist(err) {
		os.Mkdir(endpointDir, os.ModePerm)
	}

	// Get current date and time
	currentTime := time.Now()
	year := currentTime.Year()
	month := currentTime.Month().String()
	day := currentTime.Day()

	// Create log directory for the year
	yearDir := filepath.Join(endpointDir, fmt.Sprintf("%d", year))
	if _, err := os.Stat(yearDir); os.IsNotExist(err) {
		os.Mkdir(yearDir, os.ModePerm)
	}

	// Create log directory for the month
	monthDir := filepath.Join(yearDir, month)
	if _, err := os.Stat(monthDir); os.IsNotExist(err) {
		os.Mkdir(monthDir, os.ModePerm)
	}

	// Create log file path
	logFileName := fmt.Sprintf("%d-%s-%d-%s.log", year, month, day, endpointName)
	logFilePath := filepath.Join(monthDir, logFileName)

	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %s", err)
		return
	}
	defer logFile.Close()

	// Create logger
	logger := log.New(logFile, "", 0)

	// OLD Log message format
	// logMessage := fmt.Sprintf(
	// 	"DATE: %d-%s-%d\n"+
	// 		"USER: %s\n"+
	// 		"RETCODE: %s\n"+
	// 		"METHOD: %s\n"+
	// 		"ENDPOINT: %s\n"+
	// 		"REQUEST BODY: %s\n"+
	// 		"RESPONSE BODY: %s\n"+
	// 		"MESSAGE: %s\n"+
	// 		"ERROR MESSAGE: %s\n"+
	// 		"%v\n"+
	// 		"- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n",
	// 	year, month, day, user, retCode, method, endpoint, requestBody, responseBody, message, messagErr, responsErr)

	req := map[string]string{}
	json.Unmarshal(requestBody, &req)
	resp := map[string]string{}
	json.Unmarshal(responseBody, &resp)

	// GET THE LOG LEVEL
	var level string
	retCodeInt, _ := strconv.Atoi(retCode)
	if retCodeInt >= 200 && retCodeInt < 300 {
		level = "INFO"
	} else if retCodeInt >= 500 {
		level = "FATAL"
	} else if retCodeInt >= 300 && retCodeInt < 400 {
		level = "ERROR"
		logrus.Error("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
		logrus.Error(retCode, " |   ", endpoint, "     |     ", responsErr.Error())
		logrus.Error("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	} else {
		level = "WARN"
	}
	//logrus
	toBeLog := LogStruct{
		Level:         level,
		Date:          fmt.Sprintf("%d-%s-%d", year, month, day),
		Username:      user,
		Retcode:       retCode,
		Method:        method,
		Endpoint:      endpoint,
		Request_body:  req,
		Response_body: resp,
		Message:       message,
		Error_message: messagErr,
		Error:         responsErr,
	}

	toBeLogByte, marshallErr := json.Marshal(toBeLog)
	if marshallErr != nil {
		logrus.Fatal(marshallErr)
	}

	// Write log message to file
	logger.Print(string(toBeLogByte))
}

type LogStruct struct {
	Level         string            `json:"level"`
	Date          string            `json:"date"`
	Username      string            `json:"username"`
	Retcode       string            `json:"retcode"`
	Method        string            `json:"method"`
	Endpoint      string            `json:"endpoint"`
	Request_body  map[string]string `json:"request_body"`
	Response_body map[string]string `json:"response_body"`
	Message       string            `json:"message"`
	Error_message string            `json:"error_message"`
	Error         error             `json:"error"`
}
