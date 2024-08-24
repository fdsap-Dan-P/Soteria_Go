package validations

import (
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
)

func APIKeyValidation(apiKey, username, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) (response.ReturnModel, response.ApplicationDetails) {
	funcName := "Validate API Key"

	appDetails := response.ApplicationDetails{}
	// check if api key has value
	if strings.TrimSpace(apiKey) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, reqBody, []byte(""), "API Key Authorization Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	// check from the database if exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE api_key = ?", apiKey).Scan(&appDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	if appDetails.App_id == 0 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, reqBody, []byte(""), "API Key Authorization Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	// marshal the response
	appDetailsByte, marshallErr := json.Marshal(appDetails)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "Marshalling App Details Failed", marshallErr, marshallErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	middleware.ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, reqBody, appDetailsByte, "Successful", "", nil)

	successResp := response.ReturnModel{Data: response.DataModel{IsSuccess: true}}

	return successResp, appDetails
}
