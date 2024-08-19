package sendotp

import (
	"fmt"
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

// @Summary    		Set OTP via Email
// @Description 	A function that will send OTP via Email for reseting password
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string	true	"User Username"
// @Success     	200  {object} response.ResponseModel
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password/send-otp/via-email [get]
func SendOTPviaEmail(c *fiber.Ctx) error {
	username := c.Params("username")
	sendOTPRequest := request.SendingOTPRequest{}
	remark := response.DBFuncResponse{}
	userInfoResponse := response.UserInfo{}
	otpResponse := response.OTPResponse{}
	otpInterval := response.SystemConfigurationResponse{}
	isBlocked := response.ReturnModel{}

	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()
	userActivity := "Request OTP Via Email"

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_info WHERE username = ?", username).Scan(&userInfoResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userInfoResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user is not blocked before requested an otp
	isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "request", "email", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, sendOTPRequest.Institution_code, sendOTPRequest.Application_code, []byte(""), false)
	if !isBlocked.Data.IsSuccess {
		return c.JSON(isBlocked)
	}

	otp, otpCreated := middleware.OTPGeneration()

	// check first if user has already data in the table
	if fetchErr := database.DBConn.Raw("SELECT * FROM reset_password_validation WHERE user_id = ?", userInfoResponse.User_id).Scan(&otpResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if otpResponse.User_id == 0 {
		// insert to the database
		if insErr := database.DBConn.Raw("SELECT add_reset_password_validation (?, ?, ?, ?, ?) AS remark", userInfoResponse.User_id, otp, "email", "", otpCreated).Scan(&remark).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", insErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	} else {
		// update user's otp
		if insErr := database.DBConn.Raw("SELECT update_user_otp(?, ?, ?, ?, ?) AS remark", otp, "email", "", otpCreated, userInfoResponse.User_id).Scan(&remark).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", insErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, []byte(""), []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	// check if user is not yet blocked after requested an otp
	isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "request", "email", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, sendOTPRequest.Institution_code, sendOTPRequest.Application_code, []byte(""), true)
	if !isBlocked.Data.IsSuccess {
		return c.JSON(isBlocked)
	}

	// get otp interval
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'otp_req_interval' AND config_insti_code = ? AND config_app_code = ?", sendOTPRequest.Institution_code, sendOTPRequest.Application_code).Scan(&otpInterval).Error; fetchErr != nil {
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

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(userInfoResponse.Username, username, userActivity, "", "", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(username, funcName, "206", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Sent OTP", "", nil)
	return c.JSON(response.ReturnModel{
		RetCode: "206",
		Message: "Successfully Sent OTP",
		Data: errors.ErrorModel{
			Minutes: remainingMinutes,
			Seconds: remainingSeconds,
		},
	})
}
