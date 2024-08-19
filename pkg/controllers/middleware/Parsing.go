package middleware

import (
	"strconv"
	"time"
)

func ParseTime(timeStr, username, funcName, methodUsed, endpoint string) time.Time {
	timeTime, parsErr := time.Parse("2006-01-02 15:04:05.999999", timeStr)
	if parsErr != nil {
		ActivityLogger(username, funcName, "301", methodUsed, endpoint, []byte(timeStr), []byte(""), "Successful", "", parsErr)
	}

	return timeTime
}

func ParseDate(timeStr, username, funcName, methodUsed, endpoint string) time.Time {
	timeTime, parsErr := time.Parse("2006-01-02", timeStr)
	if parsErr != nil {
		ActivityLogger(username, funcName, "301", methodUsed, endpoint, []byte(timeStr), []byte(""), "Successful", "", parsErr)
	}

	return timeTime
}

func StrToFloat64(text, username, funcName, methodUsed, endpoint string) float64 {
	floatValue, parsErr := strconv.ParseFloat(text, 64)
	if parsErr != nil {
		ActivityLogger(username, funcName, "301", methodUsed, endpoint, []byte(text), []byte(""), "Successful", "", parsErr)
	}
	return floatValue
}
