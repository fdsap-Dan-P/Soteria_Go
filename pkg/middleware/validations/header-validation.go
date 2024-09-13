package validations

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
)

func HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint string) (response.ReturnModel, response.HeaderValidationResponse) {
	validationResponse := response.HeaderValidationResponse{}
	userTokenDetails := response.UserTokenDetails{}
	instiDetails := response.InstitutionDetails{}

	token := strings.TrimPrefix(authHeader, "Bearer")
	tokenString := strings.TrimSpace(token)

	fmt.Println("Auth header: ", authHeader)
	fmt.Println("Token string: ", tokenString)

	if strings.TrimSpace(authHeader) == "" || tokenString == "" {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "111", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	// validate token
	tokenValidatedStatus := middleware.ParseToken(tokenString, "", moduleName, methodUsed, endpoint)
	if !tokenValidatedStatus.Data.IsSuccess {
		return tokenValidatedStatus, validationResponse
	}

	apiKeyValidatedStatus, appDetails := APIKeyValidation(apiKey, tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, "", funcName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return apiKeyValidatedStatus, validationResponse
	}

	// check if token was stored
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM logs.user_tokens WHERE (username = ? OR staff_id = ?) AND token = ? AND insti_code = ? AND app_code = ?", tokenValidatedStatus.Message, tokenValidatedStatus.Message, tokenString, tokenValidatedStatus.Data.Message, appDetails.Application_code).Scan(&userTokenDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	if userTokenDetails.Token_id == 0 {
		returnMessage := middleware.ResponseData(tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, appDetails.Application_code, moduleName, funcName, "116", methodUsed, endpoint, []byte(""), []byte(""), "Terminated Token", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	// Get Branch Details
	if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", tokenValidatedStatus.Data.Message).Scan(&instiDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	if instiDetails.Institution_id == 0 {
		returnMessage := middleware.ResponseData(tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "Institution Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	// set up the response
	validationResponse = response.HeaderValidationResponse{
		Username:   tokenValidatedStatus.Message,
		Insti_code: instiDetails.Institution_code,
		Insti_name: instiDetails.Institution_name,
		App_code:   appDetails.Application_code,
		App_name:   appDetails.Application_name,
	}

	// marshal the response
	validationResponseByte, marshalErr := json.Marshal(validationResponse)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(tokenValidatedStatus.Message, tokenValidatedStatus.Data.Message, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Response Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, validationResponse
		}
	}

	// log the header validation
	middleware.ActivityLogger(validationResponse.Username, validationResponse.Insti_code, validationResponse.App_code, moduleName, funcName, "215", methodUsed, endpoint, []byte(""), validationResponseByte, "Successfully Validated Headers", "", nil)

	successResp := response.ReturnModel{Data: response.DataModel{IsSuccess: true}}

	return successResp, validationResponse
}
