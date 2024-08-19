package sendotp

import (
	"encoding/json"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

// @Summary    		Get OTP Channels
// @Description 	A function that will provide the list of available otp channels
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string	true	"User Username"
// @Success     	200  {object} response.SystemConfigurationResponse
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password/method [get]
func OTPMethodList(c *fiber.Ctx) error {
	username := c.Params("username")
	OTPMethodRequest := request.SendingOTPRequest{}
	userAccountResponse := response.UserAccountResponse{}
	configResponse := []response.SystemConfigurationResponse{}

	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()

	// check if user exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE username = ?", username).Scan(&userAccountResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userAccountResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get the method of reseting password
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameter.system_config WHERE config_code LIKE 'via%' AND config_value = 'true' AND config_insti_code = ? AND config_app_code = ?", OTPMethodRequest.Institution_code, OTPMethodRequest.Application_code).Order("config_name").Scan(&configResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	otpMethodOptions := make([]string, len(configResponse))
	for i, b := range configResponse {
		otpMethodOptions[i] = b.Config_name
	}

	// marshall the response body
	configResponseByte, marshallErr := json.Marshal(otpMethodOptions)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Response Failed", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	middleware.ActivityLogger(username, "Reset Password", "200", methodUsed, endpoint, (configResponseByte), configResponseByte, "Successful", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data:    configResponse,
	})
}
