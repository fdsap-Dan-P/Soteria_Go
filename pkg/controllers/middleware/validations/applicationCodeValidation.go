package validations

import (
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
)

func ValidateAppCode(appCode, username, funcName, methodUsed, endpoint string) response.ReturnModel {
	appResponse := response.ApplicationResponse{}
	successResponse := response.ReturnModel{}

	if strings.TrimSpace(appCode) == "" || strings.TrimSpace(appCode) == "null" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "Application Code Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// check if application code exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE app_code = ?", appCode).Scan(&appResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if appResponse.App_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "Application Code Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	successResponse.Data.IsSuccess = true
	return (successResponse)
}
