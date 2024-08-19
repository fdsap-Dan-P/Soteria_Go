package validations

import (
	"fmt"
	"regexp"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"
)

func ValidatePassword(password, instiCode, appCode string) bool {
	// Fetch configurations from the database
	configResponses := []response.SystemConfigurationResponse{}
	database.DBConn.Raw("SELECT * FROM paremeters.system_config WHERE config_value = 'true' AND config_code IN('pass_lowc','pass_upc','pass_num','pass_sym') AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&configResponses)
	fmt.Println("response :", configResponses)

	// Regular expressions to check for lowercase, uppercase, symbol, and number
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	symbolRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	numberRegex := regexp.MustCompile(`[0-9]`)

	// Check each configuration against the password
	for _, configResponse := range configResponses {
		switch configResponse.Config_code {
		case "pass_lowc":
			fmt.Println("LOWER LIST")
			if !lowercaseRegex.MatchString(password) {
				fmt.Println("LOWER FALSE")
				return false
			}
		case "pass_upc":
			fmt.Println("UPPER LIST")
			if !uppercaseRegex.MatchString(password) {
				fmt.Println("UPPER FALSE")
				return false
			}
		case "pass_sym":
			fmt.Println("SYMBOL LIST")
			if !symbolRegex.MatchString(password) {
				fmt.Println("SYMBOL FALSE")
				return false
			}
		case "pass_num":
			fmt.Println("NUMM LIST")
			if !numberRegex.MatchString(password) {
				fmt.Println("NUMM FALSE")
				return false
			}
		}
	}

	return true
}

func IsPasswordExpired(username, funcName, methodUsed, endpoint, currentDateTime, Last_password_reset, instiCode, appCode string, reqByte []byte) response.ReturnModel {
	passwordDuration := response.SystemConfigurationResponse{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'pass_exp' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&passwordDuration).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	passwordLastChange := middleware.ParseTime(Last_password_reset, username, funcName, methodUsed, endpoint)
	dateNowTimeTime := middleware.ParseTime(currentDateTime, username, funcName, methodUsed, endpoint)

	if strings.TrimSpace(passwordDuration.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (reqByte), []byte(""), "Password Expiration Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	passDuration, parsErr := strconv.Atoi(passwordDuration.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, (reqByte), []byte(""), "Parsing Password Expiration Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	// Convert days to time.Duration
	duration := time.Duration(passDuration) * 24 * time.Hour

	// Calculate the expiration time by adding the duration to the lastChange time
	expirationTime := passwordLastChange.Add(duration)

	// Check if the current date and time is after the expiration time
	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:   "Is Password Expired?",
			IsSuccess: dateNowTimeTime.After(expirationTime), // password expired?
			Error:     nil,
		},
	}
	// return dateNowTimeTime.After(expirationTime),
}

// ----------- FOR DAYS ----------------//
func PasswordExpireChecker(durationInDays int, lastChange, dateNow time.Time) bool {
	// Convert days to time.Duration
	duration := time.Duration(durationInDays) * 24 * time.Hour

	// Calculate the expiration time by adding the duration to the lastChange time
	expirationTime := lastChange.Add(duration)

	// Check if the current date and time is after the expiration time
	return dateNow.After(expirationTime)
}

// ------------ TESTING MINUTES ---------------//
// func IsPasswordExpired(durationInMinutes int, lastChange, dateNow time.Time) bool {
// 	// Convert minutes to time.Duration
// 	duration := time.Duration(durationInMinutes) * time.Minute

// 	// Calculate the expiration time by adding the duration to the lastChange time
// 	expirationTime := lastChange.Add(duration)

// 	fmt.Println("lastChange: ", lastChange)
// 	fmt.Println("expirationTime: ", expirationTime)
// 	fmt.Println("currentDateTime: ", dateNow)
// 	fmt.Println("Is expired: ", dateNow.After(expirationTime))

// 	// Check if the current date and time is after the expiration time
// 	return dateNow.After(expirationTime)
// }

func PasswordMinCharacter(password, instiCode, appCode string) (bool, error) {
	configResponse := response.SystemConfigurationResponse{}
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'pass_min' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&configResponse).Error; fetchErr != nil {
		return false, fetchErr
	}

	minCharacter, parsErr := strconv.Atoi(configResponse.Config_value)
	if parsErr != nil {
		return false, parsErr
	}

	if len(password) < minCharacter {
		return false, nil
	}

	return true, nil
}
