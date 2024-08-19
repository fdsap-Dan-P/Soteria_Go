package sendotp

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/changePassword/aggelos"
	"soteria_go/pkg/controllers/changePassword/restrictions"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// @Summary    		Set OTP via SMS
// @Description 	A function that will send OTP via SMS for reseting password
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string	true	"User Username"
// @Success     	200  {object} response.ResponseModel
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password/send-otp/via-sms [get]
func SendOTPviaSMS(c *fiber.Ctx) error {
	username := c.Params("username")
	sendOTPRequest := request.SendingOTPRequest{}

	userInfoResponse := response.UserInfo{}
	otpInterval := response.SystemConfigurationResponse{}
	isBlocked := response.ReturnModel{}

	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()

	userActivity := "Request OTP via SMS"

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_info WHERE username = ?", username).Scan(&userInfoResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userInfoResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", fmt.Errorf(username))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user is not yet restricted  before requesting otp
	isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "request", "sms", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, sendOTPRequest.Institution_code, sendOTPRequest.Application_code, []byte(""), false)

	// Login to the first endpoint and get the token
	aggelosLog := aggelos.AggelosLogin(username, funcName, methodUsed, endpoint)
	if !aggelosLog.Data.IsSuccess {
		return c.JSON(aggelosLog)
	}

	// Use the token to send OTP
	sendOTPreq := aggelos.SendOTP(aggelosLog.Data.Message, userInfoResponse.Phone_no, username, funcName, methodUsed, endpoint)
	if !sendOTPreq.Data.IsSuccess {
		return c.JSON(sendOTPreq)
	}

	// check if user not blocked after requested otp
	isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "request", "sms", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, sendOTPRequest.Institution_code, sendOTPRequest.Application_code, []byte(""), true)
	if !isBlocked.Data.IsSuccess {
		return c.JSON(isBlocked)
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(userInfoResponse.Username, username, userActivity, "", "", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	sendOTPrespByte, marshallErr := json.Marshal(sendOTPreq)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get otp interval
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'otp_req_interval' AND config_institution_code = ? AND config_app_code = ?", sendOTPRequest.Institution_code, sendOTPRequest.Application_code).Scan(&otpInterval).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(otpInterval.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "121", methodUsed, endpoint, []byte(""), []byte(""), "No Data Available", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	otpIntervalValue, parsErr := strconv.Atoi(otpInterval.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing OTP Interval Value Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	remainingMinutes := otpIntervalValue / 60
	remainingSeconds := otpIntervalValue % 60

	middleware.ActivityLogger(username, funcName, sendOTPreq.RetCode, methodUsed, endpoint, []byte(""), sendOTPrespByte, sendOTPreq.Message, "", nil)
	return c.JSON(response.ReturnModel{
		RetCode: sendOTPreq.RetCode,
		Message: sendOTPreq.Message,
		Data: errors.ErrorModel{
			Minutes: remainingMinutes,
			Seconds: remainingSeconds,
		},
	})
}
