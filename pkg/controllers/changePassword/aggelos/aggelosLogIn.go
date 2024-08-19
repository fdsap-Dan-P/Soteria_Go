package aggelos

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
)

func AggelosLogin(username, funcName, methodUsed, endpoint string) response.ReturnModel {
	url := "https://dev-mercury.fortress-asya.com:17000/api/v1/user/auth"
	requestBody := map[string]string{
		"username": "cagabay_report",
		"password": "B@kA1_r3p0rT",
	}
	requestBodyBytes, marshallErr := json.Marshal(requestBody)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	resp, respErr := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if respErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", respErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		returnMessage := middleware.ResponseData(username, funcName, "317", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		returnMessage := middleware.ResponseData(username, funcName, "310", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	token := tokenResponse["token"].(string)
	middleware.ActivityLogger(username, funcName, "201", methodUsed, endpoint, (requestBodyBytes), []byte(""), "Successfully Logged In to Aggelos", "", nil)
	return response.ReturnModel{
		RetCode: "201",
		Message: "Successfully Logged In to Aggelos",
		Data: errors.ErrorModel{
			Message:   token,
			IsSuccess: true,
			Error:     nil,
		},
	}
}
