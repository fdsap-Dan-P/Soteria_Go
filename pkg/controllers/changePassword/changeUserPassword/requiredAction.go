package changeuserpassword

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/controllers/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// @Summary    		Changing User's Password
// @Description 	A function that will allow the users to change their current password
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string						true	"User Username"
// @Param			Body		body  request.ResetPasswordRequest	true	"Request Body"
// @Success     	200  {object} response.ResponseModel
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password [post]
func ChangePassword(c *fiber.Ctx) error {
	username := c.Params("username")
	remark := response.DBFuncResponse{}
	newPasswordRequest := request.ResetPasswordRequest{}
	otpResponse := response.OTPResponse{}
	userAccountResponse := response.UserAccountResponse{}
	lastPasswordUsed := []response.LastUsed{}
	configResponse := response.SystemConfigurationResponse{}
	userProfileResponse := response.UserProfilesResponse{}
	currentDateTime := middleware.GetDateTime()
	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()

	if !currentDateTime.Data.IsSuccess {
		return c.JSON(currentDateTime)
	}

	if parsErr := c.BodyParser(&newPasswordRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing New Password Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshall the request body
	newHashedPasswordByte, marshallErr := json.Marshal(newPasswordRequest)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Failed", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newPasswordRequest.Otp) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "OTP Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newPasswordRequest.Password) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "New Password Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newPasswordRequest.Password) != "" && strings.TrimSpace(newPasswordRequest.Confirm_password) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "Confirm Password Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if passwords are matched
	if newPasswordRequest.Password != newPasswordRequest.Confirm_password {
		returnMessage := middleware.ResponseData(username, funcName, "119", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check password length
	isPasswordAboveMin, fetchErr := validations.PasswordMinCharacter(newPasswordRequest.Password, newPasswordRequest.Institution_code, newPasswordRequest.Application_code)
	if fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if !isPasswordAboveMin {
		returnMessage := middleware.ResponseData(username, funcName, "103", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate inputted password
	isPassValid := validations.ValidatePassword(newPasswordRequest.Password, newPasswordRequest.Institution_code, newPasswordRequest.Application_code)
	if !isPassValid {
		returnMessage := middleware.ResponseData(username, funcName, "103", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE username = ?", username).Scan(&userAccountResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userAccountResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "User Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Check if otp is validated
	if fetchErr := database.DBConn.Raw("SELECT * FROM reset_password_validation WHERE user_id = ? AND otp = ?", userAccountResponse.User_id, newPasswordRequest.Otp).Scan(&otpResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if otpResponse.Status != "Valid OTP!" {
		returnMessage := middleware.ResponseData(username, funcName, "405", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fmt.Errorf("invalid otp"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check max number of password reusability
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'pass_reuse' AND institution_code = ? AND application_code = ?", newPasswordRequest.Institution_code, newPasswordRequest.Application_code).Scan(&configResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// converter config value to numeric
	maxPasswordReuse, parsErr := strconv.Atoi(configResponse.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "Parsing Password Limit Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// New Pass must check reusability
	if fetchErr := database.DBConn.Raw("SELECT * FROM password_reuse(?, ?)", userAccountResponse.User_id, maxPasswordReuse).Scan(&lastPasswordUsed).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	newHashedPassword := hash.SHA256(newPasswordRequest.Password)

	// check if new password matches to previous password
	for _, previouseUserPassword := range lastPasswordUsed {
		if newHashedPassword == previouseUserPassword.Password_hash {
			returnMessage := middleware.ResponseData(username, funcName, "113", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	if insErr := database.DBConn.Raw("SELECT update_user_password (?, ?, ?, ?) AS remark", userAccountResponse.User_id, newHashedPassword, false, currentDateTime.Data.Message).Scan(&remark).Error; insErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", insErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user personal info
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_profiles WHERE user_id = ?", userAccountResponse.User_id).Scan(&userProfileResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// delete all bad attempts
	if deletErr := database.DBConn.Raw("SELECT delete_from_login_logs(?, ?) AS remark", userAccountResponse.User_id, false).Scan(&remark).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", deletErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// delete being restricted
	if deletErr := database.DBConn.Raw("SELECT delete_from_restricted(?, ?) AS remark", userAccountResponse.User_id, "login").Scan(&remark).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", deletErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	middleware.ActivityLogger(username, funcName, "204", methodUsed, endpoint, (newHashedPasswordByte), []byte(""), "Password Successfully Updated", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: "204", // since updating password, its 204
		Message: "Successfully Updated",
	})
}
