package loginmanagement

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/controllers/middleware/validations"

	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MainUserLogIn(c *fiber.Ctx) error {
	userCredential := request.LogInRequest{}
	userAccount := response.UserAccountResponse{}
	userPassword := response.UserPasswordsResponse{}
	userStatus := response.UserStatusResponse{}
	uponLoginInfo := response.LogInResponse{}
	soteriaResponse := response.LogInSuccess{}

	currentDateTime := middleware.GetDateTime().Data.Message

	methodUsed := c.Method()
	endpoint := c.Path()
	userIP := c.IP()
	funcName := "User Log"

	userActivity := "Logging In"

	// get user credentials
	if parsErr := c.BodyParser(&userCredential); parsErr != nil {
		returnMessage := middleware.ResponseData(userCredential.Username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing User Credential Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshall the request body
	userCredentialByte, marshallErr := json.Marshal(userCredential)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(userCredential.Password, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Failed", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check the insti code
	isInstiCodeValid := validations.ValidateInstiCode(userCredential.Institution_code, userCredential.Username, funcName, methodUsed, endpoint)
	if !isInstiCodeValid.Data.IsSuccess {
		return c.JSON(isInstiCodeValid)
	}

	isAppCodeValid := validations.ValidateAppCode(userCredential.Application_code, userCredential.Username, funcName, methodUsed, endpoint)
	if !isAppCodeValid.Data.IsSuccess {
		return c.JSON(isAppCodeValid)
	}

	// check if username is filled
	if strings.TrimSpace(userCredential.Username) == "" {
		returnMessage := middleware.ResponseData(userCredential.Password, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "Username or Employee ID Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if password is filled
	if strings.TrimSpace(userCredential.Password) == "" {
		returnMessage := middleware.ResponseData(userCredential.Password, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "Password Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if username exist and get user id
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE username = ? OR staff_id = ? OR email = ? OR phone_no = ?", userCredential.Username, userCredential.Username, userCredential.Username, userCredential.Username).Scan(&userAccount).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "302", methodUsed, endpoint, (userCredentialByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// if username does not exist terminate
	if userAccount.User_id == 0 {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "404", methodUsed, endpoint, (userCredentialByte), []byte(""), "User Not Found", fmt.Errorf("invalid username"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user status
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_status WHERE status_id = ?", userAccount.Status_id).Scan(&userStatus).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "302", methodUsed, endpoint, (userCredentialByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user status exist in db
	if userStatus.Status_id == 0 {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "404", methodUsed, endpoint, (userCredentialByte), []byte(""), "User Status Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user account is active else return
	if userStatus.Status_name != "Active" {
		errMessage := fmt.Sprintf("%s User Account", userStatus.Status_name)
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "123", methodUsed, endpoint, (userCredentialByte), []byte(""), errMessage, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user is not currently logged in
	if userAccount.Is_active {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "105", methodUsed, endpoint, (userCredentialByte), []byte(""), "", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// hashed the user inputted password
	hashInputtedPassword := hash.SHA256(userCredential.Password)

	// get user password details
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_passwords WHERE user_id = ? ORDER BY created_at DESC", userAccount.User_id).Scan(&userPassword).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "302", methodUsed, endpoint, (userCredentialByte), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user has registered password
	if userPassword.User_id == 0 {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "404", methodUsed, endpoint, (userCredentialByte), []byte(""), "User Not Found", fmt.Errorf("invalid password"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if hashInputtedPassword != userPassword.Password_hash {
		if strings.TrimSpace(userAccount.Last_login) == "" && strings.TrimSpace(userPassword.Last_password_reset) == "" && userPassword.Requires_password_reset {
			middleware.ActivityLogger(userAccount.Username, funcName, "126", methodUsed, endpoint, userCredentialByte, []byte(""), "Validation Failed", "New User Account", nil)

			sessionId := middleware.CreateAndStoreSessionId(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userIP, userAccount.User_id, userCredentialByte)
			if !sessionId.Data.IsSuccess {
				return c.JSON(sessionId)
			}
			return c.JSON(response.ReturnModel{
				RetCode: "126",
				Message: "Validation Failed",
				Data: errors.ErrorModel{
					Message:    "New User Account",
					User_id:    userAccount.User_id,
					Session_id: sessionId.Data.Message,
					Is_active:  true, // session id status
				},
			})
		} // new user inputted invalid temporary credential

		if strings.TrimSpace(userPassword.Last_password_reset) == "" {
			middleware.ActivityLogger(userAccount.Username, funcName, "126", methodUsed, endpoint, userCredentialByte, []byte(""), "Validation Failed", "Unlocked User Account", nil)

			sessionId := middleware.CreateAndStoreSessionId(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userIP, userAccount.User_id, userCredentialByte)
			if !sessionId.Data.IsSuccess {
				return c.JSON(sessionId)
			}
			return c.JSON(response.ReturnModel{
				RetCode: "126",
				Message: "Validation Failed",
				Data: errors.ErrorModel{
					Message:    "Unlocked User Account",
					User_id:    userAccount.User_id,
					Session_id: sessionId.Data.Message,
					Is_active:  true, // session id status
				},
			})
		} // unlocked, required reset password

		// --- Validatioin Here --- //
		isLocked := ValidateUserLoginAttempt(userAccount.Username, userActivity, funcName, methodUsed, endpoint, currentDateTime, userIP, userCredential.Institution_code, userCredential.Application_code, userAccount.User_id, userCredentialByte)
		if !isLocked.Data.IsSuccess {
			return c.JSON(isLocked)
		}
		// --- Validatio n Ends --- //
	}

	// check if user is new user, current unlock or user's password expired
	if strings.TrimSpace(userAccount.Last_login) == "" && strings.TrimSpace(userPassword.Last_password_reset) == "" && userPassword.Requires_password_reset {
		middleware.ActivityLogger(userAccount.Username, funcName, "101", methodUsed, endpoint, userCredentialByte, []byte(""), "Validation Failed", "Password Expired", nil)

		sessionId := middleware.CreateAndStoreSessionId(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userIP, userAccount.User_id, userCredentialByte)
		if !sessionId.Data.IsSuccess {
			return c.JSON(sessionId)
		}
		return c.JSON(response.ReturnModel{
			RetCode: "101",
			Message: "Validation Failed",
			Data: errors.ErrorModel{
				Message:    "Unlocked User Account",
				User_id:    userAccount.User_id,
				Session_id: sessionId.Data.Message,
				Is_active:  true, // session id status
			},
		})
	} // first login

	if strings.TrimSpace(userAccount.Last_login) != "" && strings.TrimSpace(userPassword.Last_password_reset) == "" && userPassword.Requires_password_reset {
		middleware.ActivityLogger(userAccount.Username, funcName, "122", methodUsed, endpoint, userCredentialByte, []byte(""), "Validation Failed", "Password Expired", nil)

		sessionId := middleware.CreateAndStoreSessionId(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userIP, userAccount.User_id, userCredentialByte)
		if !sessionId.Data.IsSuccess {
			return c.JSON(sessionId)
		}
		return c.JSON(response.ReturnModel{
			RetCode: "122",
			Message: "Validation Failed",
			Data: errors.ErrorModel{
				Message:    "Unlocked User Account",
				User_id:    userAccount.User_id,
				Session_id: sessionId.Data.Message,
				Is_active:  true, // session id status
			},
		})
	} // unlocked, correct inputted temporary credential, required reset password

	// check if password is expired
	isPasswordExpired := validations.IsPasswordExpired(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userPassword.Last_password_reset, userCredential.Institution_code, userCredential.Application_code, userCredentialByte)
	if isPasswordExpired.Data.IsSuccess && isPasswordExpired.RetCode == "200" {
		middleware.ActivityLogger(userAccount.Username, funcName, "102", methodUsed, endpoint, userCredentialByte, []byte(""), "Validation Failed", "Password Expired", nil)

		sessionId := middleware.CreateAndStoreSessionId(userAccount.Username, funcName, methodUsed, endpoint, currentDateTime, userIP, userAccount.User_id, userCredentialByte)
		if !sessionId.Data.IsSuccess {
			return c.JSON(sessionId)
		}
		return c.JSON(response.ReturnModel{
			RetCode: "102",
			Message: "Validation Failed",
			Data: errors.ErrorModel{
				Message:    "Password Expired",
				User_id:    userAccount.User_id,
				Session_id: sessionId.Data.Message,
				Is_active:  true, // session id status
			},
		})
	}

	// reset bad login attempt
	resetBadLoginAttempts := UserLogIfCorrectCredential(userAccount.Username, funcName, methodUsed, endpoint, userIP, currentDateTime, userAccount.User_id, userCredentialByte)
	if !resetBadLoginAttempts.Data.IsSuccess {
		return c.JSON(resetBadLoginAttempts)
	}

	// generate token
	token, tokenErr := middleware.GenerateToken(userAccount.Username, userCredential.Institution_code, userCredential.Application_code)
	if tokenErr != nil {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "305", methodUsed, endpoint, (userCredentialByte), []byte(""), "", tokenErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	soteriaResponse.User_id = userAccount.User_id
	soteriaResponse.Username = userAccount.Username
	soteriaResponse.Staff_id = userAccount.Staff_id
	soteriaResponse.Token = token

	// marshall the response body
	soteriaResponseByte, marshallErr := json.Marshal(soteriaResponse)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(userAccount.Username, funcName, "311", methodUsed, endpoint, (userCredentialByte), []byte(""), "Marshalling Response Failed", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(soteriaResponse.Username, soteriaResponse.Username, userActivity, "Logged Out", "Logged In", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(uponLoginInfo.User_details.Username, "User Log", "201", methodUsed, endpoint, (userCredentialByte), (soteriaResponseByte), "Successfully Logged In", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: "201",
		Message: "Successfully Logged In",
		Data:    soteriaResponse,
	})
}
