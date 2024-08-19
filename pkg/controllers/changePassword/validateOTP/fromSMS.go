package validateotp

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/changePassword/aggelos"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"

	"strings"
)

// @Summary    		Validate OTP From SMS
// @Description 	A function that will validate OTP from SMS for reseting password
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string				true	"User Username"
// @Param			Body		body  request.OTPRequest	true	"Request Body"
// @Success     	200  {object} response.ResponseModel
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password/validate-otp/via-sms [post]
func ValidateOTPfromSMS(c *fiber.Ctx) error {
	username := c.Params("username")
	userInfo := response.UserInfo{}
	otpRequest := request.OTPRequest{}
	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()

	userActivity := "Validate OTP"

	// Login to the first endpoint and get the token
	aggelosLog := aggelos.AggelosLogin(username, funcName, methodUsed, endpoint)
	if !aggelosLog.Data.IsSuccess {
		return c.JSON(aggelosLog)
	}

	if parsErr := c.BodyParser(&otpRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing OTP Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	otpByte, marshallErr := json.Marshal(otpRequest)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, (otpByte), []byte(""), "Marshalling Request Failed", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user phone number
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_info WHERE username = ?", username).Scan(&userInfo).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, (otpByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userInfo.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, (otpByte), []byte(""), "User Not Found", fmt.Errorf(username))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(otpRequest.Otp) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, (otpByte), []byte(""), "OTP Input Missing", fmt.Errorf(username))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Use the token to validate OTP
	validatErr := aggelos.ValidateOTP(username, aggelosLog.Data.Message, userInfo.Phone_no, otpRequest.Otp, funcName, methodUsed, endpoint, userActivity, userInfo.Status, otpRequest.Institution_code, otpRequest.Application_code, (otpByte), userInfo.User_id)
	if !validatErr.Data.IsSuccess {
		return c.JSON(validatErr)
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "", "", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(username, funcName, "207", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Validated OTP", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: validatErr.RetCode,
		Message: validatErr.Message,
	})
}
