package registernewuser

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/user-management/memberVerification"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NonStaffRegistraion(c *fiber.Ctx) error {
	userInstiCode := ""
	newUserRequest := request.UserRegistrationRequest{}
	userDetailValidation := response.UserDetails{}
	userDetails := response.UserDetails{}
	instiDetails := response.InstitutionDetails{}

	remark := response.DBFuncResponse{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Register Non-Staff User"

	// extract headers
	apiKey := c.Get("X-API-Key")

	// validate the api key
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return c.JSON(apiKeyValidatedStatus)
	}

	// parse the request body
	if parsErr := c.BodyParser(&newUserRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", appDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	newUserRequestByte, marshalErr := json.Marshal(newUserRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Username) == "" {
		returnMessage := middleware.ResponseData("", "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.First_name) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "First Name Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Last_name) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Last Name Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	isPhoneNoFormatted := response.ReturnModel{}
	if strings.TrimSpace(newUserRequest.Phone_no) != "" {
		// format the phone number
		isPhoneNoFormatted = middleware.NormalizePhoneNumber(newUserRequest.Phone_no, newUserRequest.Username, "", appDetails.Application_code, funcName, methodUsed, endpoint)
		if !isPhoneNoFormatted.Data.IsSuccess {
			return c.JSON(isPhoneNoFormatted)
		}
	}

	// format the birthdate
	isBdateFormatted := response.ReturnModel{}
	if strings.TrimSpace(newUserRequest.Birthdate) != "" {
		isBdateFormatted = middleware.FormatingDate(newUserRequest.Birthdate, newUserRequest.Username, "", appDetails.Application_code, funcName, methodUsed, endpoint)
		if !isBdateFormatted.Data.IsSuccess {
			return c.JSON(isBdateFormatted)
		}
	}

	// validate email address
	if strings.TrimSpace(newUserRequest.Email) != "" { // some project don't require email
		isEmailAddrValid := middleware.ValidateEmail(newUserRequest.Email)
		if !isEmailAddrValid {
			returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "109", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	// validate if username already exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ? AND application_code = ?", newUserRequest.Username, appDetails.Application_code).Scan(&userDetailValidation).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetailValidation.User_id != 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if phone number already exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE phone_no = ? AND application_code = ?", isPhoneNoFormatted.Data.Message, appDetails.Application_code).Scan(&userDetailValidation).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetailValidation.User_id != 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Phone Number Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Email) != "" {
		// validate if email address already exist
		if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE email = ? AND application_code = ?", newUserRequest.Email, appDetails.Application_code).Scan(&userDetailValidation).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if userDetailValidation.User_id != 0 {
			returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Email Address Already Exists", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	dmMemberVerifyReqBody := request.MemberVerificationRequest{
		First_name: newUserRequest.First_name,
		Last_name:  newUserRequest.Last_name,
		Phone_no:   isBdateFormatted.Data.Message,
		Birthdate:  isBdateFormatted.Data.Message,
	}

	// generate user's temp password
	tempPassword := middleware.PasswordGeneration()
	hashTempPassword := hash.SHA256(tempPassword)

	if strings.TrimSpace(newUserRequest.Institution_code) == "" {
		// get the CID and institution
		memberVerifyStatus, memberVerifyDetails := memberVerification.VerifyMemberFromDataMart(newUserRequest.Username, appDetails.Application_code, moduleName, methodUsed, endpoint, apiKey, dmMemberVerifyReqBody)
		if !memberVerifyStatus.Data.IsSuccess {
			return c.JSON(memberVerifyStatus)
		}

		// identify if user member or not
		//----- CODE HERE -----//

		if strings.TrimSpace(memberVerifyDetails.Data.Details.Insti_code) == "" {
			userInstiCode = "0000"
		}
	} else {
		userInstiCode = newUserRequest.Institution_code
	}

	// get the institution details
	if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", userInstiCode).Scan(&instiDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData("", "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if instiDetails.Institution_id == 0 {
		returnMessage := middleware.ResponseData("", "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// register the user
	if insertErr := database.DBConn.Raw("SELECT public.register_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, newUserRequest.First_name, newUserRequest.Middle_name, newUserRequest.Last_name, newUserRequest.Email, isPhoneNoFormatted.Data.Message, "", instiDetails.Institution_id, hashTempPassword, true, "", isBdateFormatted.Data.Message, instiDetails.Institution_code, appDetails.Application_code, appDetails.Application_id).Scan(&remark).Error; insertErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, newUserRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(newUserRequest.Username, newUserRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fmt.Errorf("%s", remark.Remark), nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user details
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ? AND application_code = ? AND institution_code = ?", newUserRequest.Username, appDetails.Application_code, userInstiCode).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// append the temp password to the user's details
	userDetails.Password = tempPassword

	// marshal the response body
	userDetailsByte, marshalErr := json.Marshal(userDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successResp := middleware.ResponseData(userDetails.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, userDetailsByte, "Successfully Registered User", nil, userDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)
}
