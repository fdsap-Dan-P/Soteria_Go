package middleware

import (
	"fmt"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"time"
)

func AuditTrailGeneration(user_username, username, activity, old_val, new_val, funcName, methodUsed, endpoint string) response.ReturnModel {
	userDetails := response.UserDetailsAuditTrail{}
	remark := response.DBFuncResponse{}

	currentDateTime := GetDateTime().Data.Message
	formattedDate, parsErr := time.Parse("2006-01-02 15:04:05.999999", currentDateTime)
	if parsErr != nil {
		returnMessage := ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Current Date Time Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	log_date := formattedDate.Format("2006-01-02")
	log_time := formattedDate.Format("15:04:05")

	if fetchErr := database.DBConn.Raw("SELECT * FROM logs.user_details_audit_trail WHERE username = ?", user_username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", fmt.Errorf(username))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT logs.add_audit_logs(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", userDetails.User_id, userDetails.Staff_id, userDetails.Username, userDetails.User_status, userDetails.Institution_name, activity, old_val, new_val, log_date, log_time).Scan(&remark).Error; fetchErr != nil {
		returnMessage := ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:   activity,
			IsSuccess: true,
			Error:     nil,
		},
	}
}
