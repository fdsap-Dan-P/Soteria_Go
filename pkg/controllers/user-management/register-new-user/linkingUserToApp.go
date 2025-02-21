package registernewuser

import (
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func LinkingUserToApp(username, instiCode, appCode, moduleName, funcName, methodUsed, endpoint string, userID, appID int, newUserRequestByte []byte) response.ReturnModel {
	userAppDetails := []response.UserApplicationDetails{}
	userAppResp := response.UserAppResponse{}

	successResp := response.ReturnModel{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_app_view WHERE username = ?", username).Scan(&userAppDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, "", appCode, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if len(userAppDetails) == 0 { // if the user is not linked to any application
		// link the user to this application
		if insErr := database.DBConn.Raw("INSERT INTO public.user_applications (user_id, application_id) VALUES (?, ?)", userID, appID).Scan(&userAppResp).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, "", appCode, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insErr, nil)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		successResp := middleware.ResponseData(username, appCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, newUserRequestByte, []byte(""), "Use Current Password in Linked Application", nil, successResp)
		return (successResp)

	} else { // check if the user is already linked to any application
		fmt.Println("TRACE 3")
		fmt.Println("LEN: ", len(userAppDetails))
		isUserLinked := false
		for _, userLinkedApp := range userAppDetails { // check if the user is already linked to this application
			if userLinkedApp.Application_code == appCode {
				isUserLinked = true
			}
		}

		if !isUserLinked { // if the user is not linked to this application
			fmt.Println("TRACE 4")
			// link the user to this application
			if insErr := database.DBConn.Raw("INSERT INTO public.user_applications (user_id, application_id) VALUES (?, ?)", userID, appID).Scan(&userAppResp).Error; insErr != nil {
				returnMessage := middleware.ResponseData(username, "", appCode, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insErr, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}

			successResp := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, []byte(""), "Use Current Password in Linked Application", nil, successResp)
			return (successResp)
		} else { // if the user is already linked to this application
			fmt.Println("TRACE 5")
			returnMessage := middleware.ResponseData(username, "", appCode, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Already Exists", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	}

	return (successResp)
}
