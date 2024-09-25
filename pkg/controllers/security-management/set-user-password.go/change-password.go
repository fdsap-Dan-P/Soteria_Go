package setuserpassword

import (
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
)

func ChangePasswordValidation(password, username, instiCode, appCode, moduleName, methodUsed, endpoint string, userId int) response.ReturnModel {
	remark := response.DBFuncResponse{}
	currentDateTime := middleware.GetDateTime().Data.Message

	funcName := "Change Password Validation"

	//  validate the new password
	isPassValid := PasswordValidation(password, instiCode, appCode, username, moduleName, methodUsed, endpoint)
	if !isPassValid.Data.IsSuccess {
		return (isPassValid)
	}

	// validate password reusability
	isPassReuse := PasswordReuseValidation(password, instiCode, appCode, username, moduleName, methodUsed, endpoint, userId)
	if !isPassReuse.Data.IsSuccess {
		return (isPassReuse)
	}

	// hash the new password
	newHashedPassword := hash.SHA256(password)

	// update the user's password
	if updatErr := database.DBConn.Raw("SELECT public.add_user_passwords(?, ?, ?, ?) AS remark", userId, newHashedPassword, false, currentDateTime).Scan(&remark).Error; updatErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", updatErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	successResp := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "203", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Updated Password", nil, nil)
	if !successResp.Data.IsSuccess {
		return (successResp)
	}

	return (successResp)
}
