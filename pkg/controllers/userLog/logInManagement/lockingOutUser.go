package loginmanagement

import (
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func ValidatingUserForLocking(user_id int, username, funcName, methodUsed, endpoint string) response.ReturnModel {
	userStatusList := []response.UserStatusResponse{}
	userAccountDetails := response.UserAccountResponse{}

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

	if updatErr := database.DBConn.Raw("UPDATE user_accounts SET status_id = ? WHERE user_id = ?", statusMap["Locked"], user_id).Scan(&userAccountDetails).Error; updatErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", updatErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	return response.ReturnModel{
		RetCode: "204",
		Message: "Successfully Updated",
		Data: errors.ErrorModel{
			Message:   "",
			IsSuccess: true,
		},
	}
}
