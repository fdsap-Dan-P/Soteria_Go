package setuserpassword

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UserInitiatedPasswordChange(c *fiber.Ctx) error {
	username := c.Params("username")
	changePasswordRequest := request.ChangePasswordRequest{}
	userDetails := response.UserDetails{}

	moduleName := "Security Management"
	funcName := "User Initiated Password Change"
	methodUsed := c.Method()
	endpoint := c.Path()

	// Extract api key from Authorization header
	apiKey := c.Get("X-API-Key")

	// validate the header
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return c.JSON(apiKeyValidatedStatus)
	}

	// check if user exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// parse the request body
	if parsErr := c.BodyParser(&changePasswordRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if new password was provided
	if strings.TrimSpace(changePasswordRequest.New_password) == "" {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "New Password Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	isPassChanged := ChangePasswordValidation(changePasswordRequest.New_password, userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint, userDetails.User_id)
	if !isPassChanged.Data.IsSuccess {
		return c.JSON(isPassChanged)
	}

	userDetails.Password = changePasswordRequest.New_password
	successResp := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Updated Password", nil, userDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)
}

func UserChangePasswordAfterExpired(c *fiber.Ctx) error {
	username := c.Params("username")
	changePasswordRequest := request.ChangePasswordRequest{}
	userDetails := response.UserDetails{}
	userPasswordDetails := response.UserPasswordDetails{}

	moduleName := "Security Management"
	funcName := "User Password Change After Expiration"
	methodUsed := c.Method()
	endpoint := c.Path()

	// Extract api key from Authorization header
	apiKey := c.Get("X-API-Key")

	// validate the header
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return c.JSON(apiKeyValidatedStatus)
	}

	// check if user exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// parse the request body
	if parsErr := c.BodyParser(&changePasswordRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshall the request body
	changePasswordRequestByte, marshalErr := json.Marshal(changePasswordRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if new password was provided
	if strings.TrimSpace(changePasswordRequest.New_password) == "" {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, changePasswordRequestByte, []byte(""), "New Password Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if old password was provided
	if strings.TrimSpace(changePasswordRequest.Old_password) == "" {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, changePasswordRequestByte, []byte(""), "Current Password Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate the old password
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_passwords WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", userDetails.User_id).Scan(&userPasswordDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, changePasswordRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userPasswordDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, changePasswordRequestByte, []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	newHashedPassword := hash.SHA256(changePasswordRequest.Old_password)

	// compare the inputted old password on whats in db
	if userPasswordDetails.Password_hash != newHashedPassword {
		returnMessage := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "103", methodUsed, endpoint, changePasswordRequestByte, []byte(""), "Invalid Current Password", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	isPassChanged := ChangePasswordValidation(changePasswordRequest.New_password, userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint, userDetails.User_id)
	if !isPassChanged.Data.IsSuccess {
		return c.JSON(isPassChanged)
	}

	userDetails.Password = changePasswordRequest.New_password
	successResp := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "204", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Updated Password", nil, userDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)
}

func ResetUserPasswordToTemporary(c *fiber.Ctx) error {
	username := c.Params("username")
	userDetails := response.UserDetails{}

	moduleName := "Security Management"
	funcName := "User Password Change After Expiration"
	methodUsed := c.Method()
	endpoint := c.Path()

	// Extract api key from Authorization header
	apiKey := c.Get("X-API-Key")

	// validate the header
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return c.JSON(apiKeyValidatedStatus)
	}

	// check if user exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	isPassSetTemp := SetTempPassword(userDetails.User_id, userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint)
	if !isPassSetTemp.Data.IsSuccess {
		return c.JSON(isPassSetTemp)
	}

	fmt.Println("- - - - - - - - PASSWORD TRACING - - - - - - - - ")
	fmt.Println("PROJECT NAME: SOTERIA")
	fmt.Println("FUNCTION NAME: ResetUserPasswordToTemporary")
	fmt.Println("TEMPOPARY PASSWORD: ", isPassSetTemp.Data.Message)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - ")

	userDetails.Password = isPassSetTemp.Data.Message
	successResp := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "204", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Updated Password", nil, userDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)

}
