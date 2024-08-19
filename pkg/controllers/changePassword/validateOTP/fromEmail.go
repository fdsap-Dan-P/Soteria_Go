package validateotp

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/changePassword/restrictions"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary    		Validate OTP From Email
// @Description 	A function that will validate OTP from Email for reseting password
// @Tags        	RESET PASSWORD
// @Accept      	json
// @Produce     	json
// @Param			username	path  string				true	"User Username"
// @Param			Body		body  request.OTPRequest	true	"Request Body"
// @Success     	200  {object} response.ResponseModel
// @Failure     	400  {object} response.ReturnModel
// @Router      	/{username}/reset-password/validate-otp/via-email [post]
func ValidateOTPfromEmail(c *fiber.Ctx) error {
	username := c.Params("username")
	remark := response.DBFuncResponse{}
	userInfoResponse := response.UserInfo{}
	otpRequest := request.OTPRequest{}
	otpResponse := response.OTPResponse{}
	currentDateTime := middleware.GetDateTime().Data.Message
	isBlocked := response.ReturnModel{}

	funcName := "Reset Password"
	methodUsed := c.Method()
	endpoint := c.Path()
	userActivity := "Validate OTP From Email"

	dateNowTime := middleware.ParseTime(currentDateTime, username, funcName, methodUsed, endpoint)

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

	if parsErr := c.BodyParser(&otpRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing User Credentials Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if otp has a value
	if strings.TrimSpace(otpRequest.Otp) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "OTP Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal request body
	requestBodyBytes, marshallErr := json.Marshal(otpRequest)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if the user is blocked before validating otp
	isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "validate", "email", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, otpRequest.Institution_code, otpRequest.Application_code, requestBodyBytes, false)
	if !isBlocked.Data.IsSuccess {
		return c.JSON(isBlocked)
	}

	// check if otp and user_id exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM reset_password_validation WHERE user_id = ? AND otp = ? AND via = 'email'", userInfoResponse.User_id, otpRequest.Otp).Scan(&otpResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, requestBodyBytes, []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if otpResponse.User_id == 0 || len(otpResponse.Otp) != 6 {
		// check if the user is blocked after validating invalid otp
		isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "validate", "email", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, otpRequest.Institution_code, otpRequest.Application_code, requestBodyBytes, true)
		if !isBlocked.Data.IsSuccess {
			return c.JSON(isBlocked)
		}

		errMsg := "Invalid OTP. You have " + isBlocked.Data.Message + " remaining attempts."
		returnMessage := middleware.ResponseData(username, funcName, "117", methodUsed, endpoint, requestBodyBytes, []byte(""), errMsg, fmt.Errorf("via email"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	otpInt, parsErr := strconv.Atoi(otpResponse.Otp)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, requestBodyBytes, []byte(""), "Parsing OTP Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	layout := "2006-01-02 15:04:05"
	creationTime, parsErr := time.Parse(layout, otpResponse.Created_at)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, requestBodyBytes, []byte(""), "Parsing Time Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Convert minutes to time.Duration
	duration := time.Duration(5) * time.Minute

	// Calculate the expiration time by adding the duration to the lastChange time
	expirationTime := creationTime.Add(duration)

	// Check if the current date and time is after the expiration time
	if dateNowTime.After(expirationTime) {
		returnMessage := middleware.ResponseData(username, funcName, "118", methodUsed, endpoint, requestBodyBytes, []byte(""), "", fmt.Errorf(fmt.Sprintf("Expiration: %v | Date Now: %v", expirationTime, currentDateTime)))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if otpInt <= 100000 || otpInt >= 999999 {
		// check if the user is blocked after validating invalid otp
		isBlocked = restrictions.ValidateUserForBlocking(userInfoResponse.User_id, "validate", "email", userActivity, userInfoResponse.Status, userInfoResponse.Username, funcName, methodUsed, endpoint, otpRequest.Institution_code, otpRequest.Application_code, requestBodyBytes, true)
		if !isBlocked.Data.IsSuccess {
			return c.JSON(isBlocked)
		}

		errMsg := "Invalid OTP. You have " + isBlocked.Data.Message + " remaining attempts."
		returnMessage := middleware.ResponseData(username, funcName, "117", methodUsed, endpoint, requestBodyBytes, []byte(""), errMsg, fmt.Errorf("otp length"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// update user's otp status
	if updatErr := database.DBConn.Debug().Raw("SELECT update_otp_status(?, ?, ?, ?, ?) AS remark", "Valid OTP!", userInfoResponse.User_id, otpRequest.Otp, "email", currentDateTime).Scan(&remark).Error; updatErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "304", methodUsed, endpoint, requestBodyBytes, []byte(""), "", updatErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, requestBodyBytes, []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "", "", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(username, funcName, "207", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Validated OTP", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: "207",
		Message: "Successfully Validated OTP",
	})
}
