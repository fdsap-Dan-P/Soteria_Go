package usermanagement

import (
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MemberVerification(c *fiber.Ctx) error {
	userRequest := request.MemberVerificationRequest{}
	userVerification := response.MemberVerificationResponse{}
	userDetails := response.UserDetails{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "HCIS User Details Provider"

	// Extraxt the headers
	apiKey := c.Get("X-API-Key")

	validationStatus, validationDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !validationStatus.Data.IsSuccess {
		return c.JSON(validationStatus)
	}

	// get the request body
	if parsErr := c.BodyParser(&userRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", validationDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fullName := "Full Name: " + userRequest.First_name + " " + userRequest.Last_name

	// marshal the request body
	userRequestByte, marshalErr := json.Marshal(userRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, userRequestByte, []byte(""), "Marshalling Request Body Failed", marshalErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.Phone_no) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "Phone Number Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.First_name) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "First Name Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.Last_name) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "Last Name Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.Birthdate) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "Birth Date Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate the phone number format
	isPhoneNumberValidated := middleware.NormalizePhoneNumber(userRequest.Phone_no, fullName, "", validationDetails.Application_code, funcName, methodUsed, endpoint)
	if !isPhoneNumberValidated.Data.IsSuccess {
		return c.JSON(isPhoneNumberValidated)
	}
	normalizedPhonenumber := isPhoneNumberValidated.Data.Message

	// format the birthdate
	isBirthDateValid := middleware.FormatingDate(userRequest.Birthdate, fullName, "", validationDetails.Application_code, moduleName, funcName, methodUsed, endpoint)
	if !isBirthDateValid.Data.IsSuccess {
		return c.JSON(isBirthDateValid)
	}
	// formattedBirthDate := isBirthDateValid.Data.Message

	// check if user is a member
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_details WHERE phone_no = ? AND first_name ILIKE ? AND last_name ILIKE ?", normalizedPhonenumber, userRequest.First_name, userRequest.Last_name).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, userRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	userVerification.Phone_no = userDetails.Phone_no
	userVerification.First_name = userDetails.First_name
	userVerification.Last_name = userDetails.Last_name
	userVerification.Birthdate = userDetails.Birthdate

	if userDetails.User_id == 0 {
		userVerification.Is_member = false
	} else {
		userVerification.Is_member = true
		userVerification.Institution_code = userDetails.Institution_code
		userVerification.Institution_name = userDetails.Institution_name
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   "Successfully Verified",
			IsSuccess: true,
			Error:     nil,
			Details:   userVerification,
		},
	})
}
