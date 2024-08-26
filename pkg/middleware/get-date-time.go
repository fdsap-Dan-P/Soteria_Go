package middleware

import (
	"soteria_go/pkg/models/response"
	"time"
)

func GetDateTime() response.ReturnModel {
	// Set the timezone to Asia/Manila
	// location, dateTimErr := time.LoadLocation("Asia/Manila")
	// if dateTimErr != nil {
	// 	return response.ReturnModel{
	// 		RetCode: "312",
	// 		Message: "Internal Server Error",
	// 		Data: errors.ErrorModel{
	// 			Message:   "Timezone Loading Failed",
	// 			IsSuccess: false,
	// 			Error:     dateTimErr,
	// 		},
	// 	}
	// }

	// Get the current time in the specified timezone
	// currentTime := time.Now().In(location)
	currentTime := time.Now()

	// Format the time using a custom format
	formattedTime := currentTime.Format("2006-01-02 15:04:05.999999")

	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   formattedTime,
			IsSuccess: true,
			Error:     nil,
		},
	}
}
