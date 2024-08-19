package aggelos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"soteria_go/pkg/controllers/middleware"

	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
)

func SendOTP(token, Phone_no, username, funcName, methodUsed, endpoint string) response.ReturnModel {
	url := "https://dev-mercury.fortress-asya.com:17000/api/v1/message/requestOtp"
	requestBody := map[string]string{
		"from":      "CARD RBI",
		"to":        Phone_no,
		"mobile":    Phone_no,
		"msg":       "Password Reset OTP \n\nYour One-Time Password is @. It expires in 5 minutes. Please do not share with anyone! \n\nData Platform Support Team",
		"appCode":   "Cagabay_RBI",
		"instiCode": "200",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		returnMessage := middleware.ResponseData(username, funcName, "317", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err != nil {
			returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", err)
			if !returnMessage.Data.IsSuccess {
				return returnMessage
			}
		}
	}

	middleware.ActivityLogger(username, funcName, "206", methodUsed, endpoint, (requestBodyBytes), []byte(""), "Successfully Sent OTP", "", nil)
	return response.ReturnModel{
		RetCode: "206",
		Message: "Successfully Sent OTP",
		Data: errors.ErrorModel{
			Message:   "",
			IsSuccess: true,
			Error:     nil,
		},
	}
}
