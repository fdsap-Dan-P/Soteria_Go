package middleware

import (
	"fmt"
	"soteria_go/pkg/models/response"
)

func NormalizePhoneNumber(phonenumber, username, instiCode, appCode, funcName, methodUsed, endpoint string) response.ReturnModel {
	moduleName := "Normalize Phone Number"
	var normalizedPhonenumber string

	if len(phonenumber) == 0 {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Phone Number Input Missing", fmt.Errorf("null phone number"), nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if phonenumber[0:1] == "0" && phonenumber[1:2] == "9" {
		normalizedPhonenumber = phonenumber
	} else if phonenumber[0:1] == "6" && len(phonenumber) == 12 {
		normalizedPhonenumber = "0" + phonenumber[2:12]
	} else if phonenumber[0:1] == "+" || phonenumber[0:1] == " " {
		normalizedPhonenumber = "0" + phonenumber[3:13]
	} else if phonenumber[0:1] == "9" {
		normalizedPhonenumber = "0" + phonenumber
	} else {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "108", methodUsed, endpoint, []byte(phonenumber), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if len(normalizedPhonenumber) != 11 {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "108", methodUsed, endpoint, []byte(phonenumber), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(phonenumber), []byte(normalizedPhonenumber), "Successful", "", nil)
	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   normalizedPhonenumber,
			IsSuccess: true,
			Error:     nil,
		},
	}
}
