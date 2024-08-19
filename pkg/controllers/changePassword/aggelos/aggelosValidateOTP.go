package aggelos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"soteria_go/pkg/controllers/changePassword/restrictions"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func ValidateOTP(username, token, phone, otp, funcName, methodUsed, endpoint, userActivity, userStatus, instiCode, appCode string, reqBody []byte, user_id int) response.ReturnModel {
	resetPasswordResponse := response.OTPResponse{}
	currentDateTime := middleware.GetDateTime().Data.Message
	remark := response.DBFuncResponse{}
	isBlocked := response.ReturnModel{}

	// check if user is blocked before validating otp
	isBlocked = restrictions.ValidateUserForBlocking(user_id, "validate", "sms", userActivity, userStatus, username, funcName, methodUsed, endpoint, instiCode, appCode, reqBody, false)
	if !isBlocked.Data.IsSuccess {
		return isBlocked
	}

	url := "https://dev-mercury.fortress-asya.com:17000/api/v1/message/validateOTP"
	requestBody := map[string]string{
		"otp":       otp,
		"mobile":    phone,
		"instiCode": "200",
	}
	requestBodyBytes, marshallErr := json.Marshal(requestBody)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if reqErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, reqBody, []byte(""), "", reqErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, respErr := client.Do(req)
	if respErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "317", methodUsed, endpoint, reqBody, []byte(""), "", respErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	defer resp.Body.Close()

	// Print response body
	responseBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "317", methodUsed, endpoint, reqBody, []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	aggelosResponse := response.ValidateOtpResponse{}
	if unmarshallErr := json.Unmarshal(responseBody, &aggelosResponse); unmarshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "310", methodUsed, endpoint, reqBody, []byte(""), "", unmarshallErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	if resp.StatusCode != http.StatusOK {
		returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, reqBody, []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	// check first if user_id already have record in the table
	if fetchErr := database.DBConn.Raw("SELECT * FROM reset_password_validation WHERE user_id = ?", user_id).Scan(&resetPasswordResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}
	if resetPasswordResponse.User_id == 0 { // insert otp and user
		if insErr := database.DBConn.Raw("SELECT add_reset_password_validation (?, ?, ?, ?, ?) AS remark", user_id, otp, "sms", aggelosResponse.Msg, currentDateTime).Scan(&remark).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", nil)
			if !returnMessage.Data.IsSuccess {
				return returnMessage
			}
		}
		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return returnMessage
			}
		}
	} else { // update user's otp
		if insErr := database.DBConn.Raw("SELECT update_user_otp(?, ?, ?, ?, ?) AS remark", otp, "sms", aggelosResponse.Msg, currentDateTime, user_id).Scan(&remark).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", nil)
			if !returnMessage.Data.IsSuccess {
				return returnMessage
			}
		}
		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return returnMessage
			}
		}
	}

	if aggelosResponse.Msg == "Expired OTP!" {
		returnMessage := middleware.ResponseData(username, funcName, "118", methodUsed, endpoint, reqBody, []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	} else if aggelosResponse.Msg == "Invalid OTP!" {
		// check if user is blocked after validating invalid otp
		isBlocked = restrictions.ValidateUserForBlocking(user_id, "validate", "sms", userActivity, userStatus, username, funcName, methodUsed, endpoint, instiCode, appCode, reqBody, true)
		if !isBlocked.Data.IsSuccess {
			return isBlocked
		}

		errMsg := "Invalid OTP. You have " + isBlocked.Data.Message + " remaining attempts."
		returnMessage := middleware.ResponseData(username, funcName, "117", methodUsed, endpoint, reqBody, []byte(""), errMsg, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage
		}
	}

	return response.ReturnModel{
		RetCode: "207",
		Message: "Successfully Validated OTP",
		Data: errors.ErrorModel{
			Message:   "Successfully",
			IsSuccess: true,
			Error:     nil,
		},
	}
}
