package logoutmanagement

import (
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func WithSessionLogOut(username, funcName, methodUsed, endpoint, sessionId string) response.ReturnModel {
	remark := response.DBFuncResponse{}
	sessionResponse := response.SessionLogs{}
	userResponse := response.UserAccountResponse{}
	currentDateTime := middleware.GetDateTime()

	userActivity := "Logged Out"
	if !currentDateTime.Data.IsSuccess {
		return (currentDateTime)
	}

	// get user id
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE username = ?", username).Scan(&userResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, " 302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// check if session is correct
	if fetchErr := database.DBConn.Raw("SELECT * FROM sessions WHERE session_id = ? AND user_id = ?", sessionId, userResponse.User_id).Scan(&sessionResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, " 302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if sessionResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "Session Id Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// make false the user's is active
	if updatErr := database.DBConn.Raw("SELECT update_concurrent_status(?, ?, ?) AS remark", false, currentDateTime.Data.Message, userResponse.User_id).Scan(&remark).Error; updatErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, " 304", methodUsed, endpoint, []byte(""), []byte(""), "", updatErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, " 304", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// delete session logs
	if deletErr := database.DBConn.Raw("SELECT delete_from_sessions (?) AS remark", userResponse.User_id).Scan(&remark).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, " 314", methodUsed, endpoint, []byte(""), []byte(""), "", deletErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, " 314", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "Logged In", "Logged Out", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return (auditTrailLogs)
	}

	middleware.ActivityLogger(username, "User Log", "202", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Log Out", "", nil)
	return (response.ReturnModel{
		RetCode: "202",
		Message: "Successfully Log Out",
		Data: errors.ErrorModel{
			IsSuccess: true,
		},
	})
}