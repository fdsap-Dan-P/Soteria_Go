package validations

import (
	"encoding/json"
	"fmt"
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/encryptDecrypt"
	"strings"
)

func APIKeyValidation(apiKey, username, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) (response.ReturnModel, response.ApplicationDetails) {
	appDetails := response.ApplicationDetails{}
	funcName := "Validate API Key"

	// check if api key has value
	if strings.TrimSpace(apiKey) == "" {
		fmt.Println("API KEY: ", apiKey)
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, reqBody, []byte(""), "API Key Authorization Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	// get the secret key
	secretKey := os.Getenv("SECRET_KEY")
	if strings.TrimSpace(secretKey) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "Secret Key Not Found in Environment", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, response.ApplicationDetails{}
		}
	}

	// encrypt the given api key
	encryptedApiKey, encryptErr := encryptDecrypt.EncryptWithSecretKey(apiKey, secretKey)
	if encryptErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "318", methodUsed, endpoint, reqBody, []byte(""), "Encrypting API Key Failed", encryptErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, response.ApplicationDetails{}
		}
	}

	// check from the database if exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE api_key = ?", encryptedApiKey).Scan(&appDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, appDetails
		}
	}

	if appDetails.Application_id == 0 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "API Key Not Found", nil, nil)
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

	fmt.Println("APP DETAILS: ", appDetails)
	return successResp, appDetails
}
