package setuserpassword

import (
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
)

func SetTempPassword(userId int, username, instiCode, appCode, moduleName, methodUsed, endpoint string) response.ReturnModel {
	remark := response.DBFuncResponse{}

	funcName := "Set User's Password to Temporary"

	// generate user's temp password
	tempPassword := middleware.PasswordGeneration()
	hashTempPassword := hash.SHA256(tempPassword)

	// register the user
	if insertErr := database.DBConn.Raw("SELECT public.add_user_passwords(?, ?, ?, ?, ?, ?) AS remark", userId, hashTempPassword, true, "", instiCode, appCode).Scan(&remark).Error; insertErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", insertErr, insertErr.Error())
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf("%s", remark.Remark), remark)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	fmt.Println("- - - - - - - - PASSWORD TRACING - - - - - - - - ")
	fmt.Println("PROJECT NAME: SOTERIA")
	fmt.Println("FUNCTION NAME: SetTempPassword")
	fmt.Println("TEMPOPARY PASSWORD: ", tempPassword)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - ")

	successResp := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "204", methodUsed, endpoint, []byte(""), []byte(""), tempPassword, nil, nil)

	return successResp
}
