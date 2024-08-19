package loginmanagement

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

func ValidateUserLoginAttempt(username, userActivity, funcName, methodUsed, endpoint, currentDateTime, userIP, instiCode, appCode string, userId int, reqBody []byte) response.ReturnModel {
	remark := response.DBFuncResponse{}
	attemptCount := response.Total{}
	lockOutAttempt := response.SystemConfigurationResponse{}
	restrictedUser := response.RestrictedResponse{}
	lockOutPeriod := response.SystemConfigurationResponse{}

	// consider as bad attempt, add to bad log
	if insErr := database.DBConn.Raw("SELECT logs.add_login_logs(?, ?, ?, ?) AS remark", userId, userIP, currentDateTime, false).Scan(&remark).Error; insErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqBody), []byte(""), "", insErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqBody), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	// count user's bad attempts
	if fetchErr := database.DBConn.Raw("SELECT COUNT(*) FROM logs.login_logs WHERE is_success = false AND user_id = ?", userId).Scan(&attemptCount).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqBody), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	// get max lockout attempt
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'lockout_attempt' AND config_insti_code = ? AND  config_app_code = ?", instiCode, appCode).Scan(&lockOutAttempt).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqBody), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if strings.TrimSpace(lockOutAttempt.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (reqBody), []byte(""), "Lockout Period Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// parse lockdown period str to int
	lockDownPeriodInt, parsErr := strconv.Atoi(lockOutAttempt.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, (reqBody), []byte(""), "Parsing Lockout Period Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// check if user reached max bad attempt
	if attemptCount.Count >= lockDownPeriodInt {
		if insErr := database.DBConn.Raw("SELECT logs.add_locked_users (?, ?, ?) AS remark", userId, "", currentDateTime).Scan(&remark).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqBody), []byte(""), "", insErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqBody), []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		if fetchErr := database.DBConn.Raw("SELECT * FROM logs.locked_users WHERE user_id = ?", userId).Scan(&restrictedUser).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqBody), []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		// start counting of lockout period
		currentTime := middleware.ParseTime(currentDateTime, username, funcName, methodUsed, endpoint)
		prevTime := middleware.ParseTime(restrictedUser.Created_at, username, funcName, methodUsed, endpoint)

		timeDiff := currentTime.Sub(prevTime)

		// get lockout period
		if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'lockout_period' AND config_insti_code = ? AND  config_app_code = ?", instiCode, appCode).Scan(&lockOutPeriod).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqBody), []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		if strings.TrimSpace(lockOutPeriod.Config_value) == "" {
			returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (reqBody), []byte(""), "Lockout Period Not Found", parsErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		// parse lockdown period str to int
		lockDownPeriodInt, parsErr := strconv.Atoi(lockOutPeriod.Config_value)
		if parsErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (reqBody), []byte(""), "Parsing Lockout Period Failed", parsErr)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		lockDownPeriodSec := lockDownPeriodInt * 60
		remainingSec := lockDownPeriodSec - int(timeDiff.Seconds())

		// Calculate remaining minutes and seconds
		remainingMinutes := remainingSec / 60
		remainingSeconds := remainingSec % 60

		if timeDiff < time.Duration(lockDownPeriodSec)*time.Second {
			// Audit Trails
			newValAuditTrail := fmt.Sprintf("Lockout Period Running: Minute[s]: %v | Second[s]: %v", remainingMinutes, remainingSeconds)
			auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "Active", newValAuditTrail, funcName, methodUsed, endpoint)
			if !auditTrailLogs.Data.IsSuccess {
				return (auditTrailLogs)
			}

			middleware.ActivityLogger(username, funcName, "114", methodUsed, endpoint, reqBody, []byte(""), "Validation Failed", fmt.Sprintf("Remaining Time: %v Minute[s] and %v Second[s]", remainingMinutes, remainingSeconds), nil)
			return (response.ReturnModel{
				RetCode: "114",
				Message: "Validation Failed",
				Data: errors.ErrorModel{
					Message: "User Lockout Period",
					Minutes: remainingMinutes,
					Seconds: remainingSeconds,
				},
			})
		} else {
			isLocked := ValidatingUserForLocking(userId, username, funcName, methodUsed, endpoint)
			if !isLocked.Data.IsSuccess {
				return (isLocked)
			}

			// Audit Trails
			auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "Active", "Locked", funcName, methodUsed, endpoint)
			if !auditTrailLogs.Data.IsSuccess {
				return (auditTrailLogs)
			}

			middleware.ActivityLogger(username, funcName, "123", methodUsed, endpoint, reqBody, []byte(""), "Validation Failed", "Locked User Account", nil)
			return (response.ReturnModel{
				RetCode: "123",
				Message: "Validation Failed",
				Data: errors.ErrorModel{
					Message:   "Locked User Account",
					IsSuccess: false,
					Error:     nil,
				},
			})
		}
	}

	returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (reqBody), []byte(""), "User Not Found 1", fmt.Errorf("invalid password"))
	if !returnMessage.Data.IsSuccess {
		return (returnMessage)
	}

	return returnMessage
}
