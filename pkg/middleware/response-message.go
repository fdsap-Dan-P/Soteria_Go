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

	fmt.Println("")
	fmt.Printf("CAGABAY UA | %s | %s | Error: %+v | %s | %s |", returnMessage.RetCode, returnMessage.Data.Message, returnMessage.Data.Error, funcName, endpoint)

	ActivityLogger(username, instiCode, appCode, moduleName, funcName, retcode, method, endpoint, []byte(reqBody), []byte(""), returnMessage.Message, returnMessage.Data.Message, error_message)
	return returnMessage
}
