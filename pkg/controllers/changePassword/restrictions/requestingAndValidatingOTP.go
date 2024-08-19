package restrictions

import (
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"
)

func ValidateUserForBlocking(userId int, otpMethod, via, userActivity, userStatus, username, funcName, methodUsed, endpoint, instiCode, appCode string, reqByte []byte, toBeInsert bool) response.ReturnModel {
	totalRequests := response.Total{}
	otpLogs := response.OTPResponse{}
	otpMaxReq := response.SystemConfigurationResponse{}
	userStatusList := []response.UserStatusResponse{}
	userAccountDetails := response.UserAccountResponse{}

	currentDateTime := middleware.GetDateTime().Data.Message
	formattedDate, parsErr := time.Parse("2006-01-02 15:04:05.999999", currentDateTime)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Current Date Time Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	dateNow := formattedDate.Format("2006-01-02")

	// is it required to insert
	if toBeInsert {
		if insErr := database.DBConn.Raw("INSERT INTO otp_request_logs(user_id, method, via, created_at) VALUES (?, ?, ?, ?)", userId, otpMethod, via, dateNow).Scan(&otpLogs).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", insErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

	}

	if otpMethod == "request" { // max 3 req combined sms and email
		// count otp request made today
		if fetchErr := database.DBConn.Raw("SELECT COUNT(*) FROM otp_request_logs WHERE user_id = ? AND method = ? AND created_at = ?", userId, otpMethod, dateNow).Scan(&totalRequests).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	} else if otpMethod == "validate" { // max 3 validation for sms and email respectively
		// count otp request made today
		if fetchErr := database.DBConn.Raw("SELECT COUNT(*) FROM otp_request_logs WHERE user_id = ? AND method = ? AND created_at = ? AND via = ?", userId, otpMethod, dateNow, via).Scan(&totalRequests).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	}

	// get the max otp request per day
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'otp_req_max' AND config_insti_code = ? AND  config_app_code = ?", instiCode, appCode).Scan(&otpMaxReq).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if strings.TrimSpace(otpMaxReq.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "121", methodUsed, endpoint, []byte(""), []byte(""), "No Maximum OTP Request Value Available", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	//convert to int
	otpMaxReqValue, parsErr := strconv.Atoi(otpMaxReq.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Maximum OTP Request Value Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// if otp requested is for a day is more than the max set
	if otpMaxReqValue <= totalRequests.Count {
		// get every id of user status
		if fetchErr := database.DBConn.Raw("SELECT * FROM user_status").Scan(&userStatusList).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		if len(userStatusList) == 0 {
			returnMessage := middleware.ResponseData(username, funcName, "121", methodUsed, endpoint, []byte(""), []byte(""), "No User Status Available", nil)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		// Create a map to store status name to status ID mapping
		statusMap := make(map[string]int)
		for _, resource := range userStatusList {
			statusMap[resource.Status_name] = resource.Status_id
		}

		if updatErr := database.DBConn.Raw("UPDATE user_accounts SET status = ? WHERE user_id = ?", statusMap["Blocked"], userId).Scan(&userAccountDetails).Error; updatErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", updatErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		returnMessage := middleware.ResponseData(username, funcName, "123", methodUsed, endpoint, []byte(""), []byte(""), "Your account has been blocked.", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	remainingAttemp := otpMaxReqValue - totalRequests.Count

	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:   fmt.Sprintf("%v", remainingAttemp),
			IsSuccess: true,
		},
	}
}
