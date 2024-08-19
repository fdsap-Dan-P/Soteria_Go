package loginmanagement

import (
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func UserLogIfCorrectCredential(username, funcName, methodUsed, endpoint, userIP, currentDateTime string, userId int, reqByte []byte) response.ReturnModel {
	remark := response.DBFuncResponse{}

	// if got correct credential, reset count of bad attempts
	if deletErr := database.DBConn.Raw("SELECT logs.delete_from_login_logs(?, ?) AS remark", userId, false).Scan(&remark).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (reqByte), []byte(""), "", deletErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (reqByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// log the users login
	if insErr := database.DBConn.Raw("SELECT logs.add_login_logs (?, ?, ?, ?) AS remark", userId, userIP, currentDateTime, true).Scan(&remark).Error; insErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqByte), []byte(""), "", insErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (reqByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:   "Reset Bad Attempt Count",
			IsSuccess: true,
			Error:     nil,
		},
	}
}
