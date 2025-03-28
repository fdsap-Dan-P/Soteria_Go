package middleware

import (
	"fmt"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
)

func ResponseData(username, instiCode, appCode, moduleName, funcName, retcode, method, endpoint string, reqBody []byte, respBsody []byte, specific_field string, error_message error, details interface{}) response.ReturnModel {
	respFromDB := response.RespFromDB{}
	returnMessage := response.ReturnModel{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.return_message WHERE retcode = ?", retcode).Scan(&respFromDB).Error; fetchErr != nil {
		return response.ReturnModel{
			RetCode: "302",
			Message: "Internal Server Error",
			Data: response.DataModel{
				Message:   "Fetching Data Failed",
				IsSuccess: false,
				Error:     fetchErr,
				Details:   fetchErr,
			},
		}
	}
	if respFromDB.Category == "" {
		return response.ReturnModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: response.DataModel{
				Message:   "Ret Code Not Found",
				IsSuccess: false,
				Error:     nil,
				Details:   nil,
			},
		}
	}

	returnMessage.RetCode = retcode
	returnMessage.Message = respFromDB.Category
	returnMessage.Data.Message = respFromDB.Error_message
	returnMessage.Data.IsSuccess = respFromDB.Is_success

	if specific_field != "" {
		returnMessage.Data.Message = specific_field
	}

	if error_message != nil {
		returnMessage.Data.Error = error_message
	}

	if details != nil {
		returnMessage.Data.Details = details
	}

	// fmt.Println()
	// fmt.Println("============== FAILED REQUEST =============")
	// fmt.Println("REF ID: ", refID)
	fmt.Println("============= REQUEST DETAILS =============")
	fmt.Println("USERNAME: ", username)
	fmt.Println("APP CODE: ", appCode)
	fmt.Println("INSTI CODE: ", instiCode)
	fmt.Println("FUNCTION NAME: ", funcName)
	// if funcName != subFuncName {
	// 	fmt.Println("SUB FUNCTION NAME: ", subFuncName)
	// }
	fmt.Println("METHOD: ", method)
	fmt.Println("ENDPOINT: ", endpoint)
	fmt.Println("REQUEST BDDY: ", string(reqBody))

	// Response Information
	fmt.Println("============ RESPONSE DETAILS =============")
	fmt.Println("RET CODE: ", returnMessage.RetCode)
	fmt.Println("RESPONSE CATEGORY: ", returnMessage.Message)
	fmt.Println("RESPONSE MESSAGE: ", returnMessage.Data.Message)
	fmt.Println("RESPONSE MESSAGE: ", string(respBsody))
	fmt.Println("ERROR MESSAGE: ", returnMessage.Data.Error)
	fmt.Println()

	ActivityLogger(username, instiCode, appCode, moduleName, funcName, retcode, method, endpoint, []byte(reqBody), []byte(""), returnMessage.Message, returnMessage.Data.Message, error_message)
	return returnMessage
}
