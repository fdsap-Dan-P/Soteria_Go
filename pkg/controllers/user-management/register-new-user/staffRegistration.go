package registernewuser

import (
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func StaffRegistration(c *fiber.Ctx) error {
	newUserRequest := request.UserRegistrationRequest{}
	userDetailValidation := response.UserDetails{}
	UserDetails := response.UserDetails{}
	instiDetails := response.InstitutionDetails{}
	remark := response.DBFuncResponse{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Register Staff User"

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

	// check if staff_id has value
	if strings.TrimSpace(newUserRequest.Staff_id) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Username) == "" {
		newUserRequest.Username = newUserRequest.Staff_id
	}

	if strings.TrimSpace(newUserRequest.Institution_code) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Code Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate the staff id format
	isStaffIdValidated := validations.StaffIdValidation(newUserRequest.Staff_id, moduleName, methodUsed, endpoint)
	if !isStaffIdValidated {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "103", methodUsed, endpoint, newUserRequestByte, []byte(""), "Invalid Employee ID", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// format the phone number
	isPhoneNoFormatted := middleware.NormalizePhoneNumber(newUserRequest.Phone_no, newUserRequest.Username, "", appDetails.Application_code, funcName, methodUsed, endpoint)
	if !isPhoneNoFormatted.Data.IsSuccess {
		return c.JSON(isPhoneNoFormatted)
	}

	// format the birthdate
	isBdateFormatted := middleware.FormatingDate(newUserRequest.Birthdate, newUserRequest.Username, "", appDetails.Application_code, funcName, methodUsed, endpoint)
	if !isBdateFormatted.Data.IsSuccess {
		return c.JSON(isBdateFormatted)
	}

	// validate email address
	isEmailAddrValid := middleware.ValidateEmail(newUserRequest.Email)
	if !isEmailAddrValid {
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "103", methodUsed, endpoint, newUserRequestByte, []byte(""), "Invalid Email Address", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if staff id already exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE staff_id = ? AND application_code = ?", newUserRequest.Staff_id, appDetails.Application_code).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetails.User_id != 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get hcis details
	hcisResponseStatus, hcisResponseDeatails := HcisInquiry(newUserRequest.Staff_id, newUserRequest.Username, "", appDetails.Application_code, "User Registration", methodUsed, endpoint, newUserRequestByte)

	// generate user's temp password
	tempPassword := middleware.PasswordGeneration()
	hashTempPassword := hash.SHA256(tempPassword)

	if hcisResponseStatus.RetCode == "405" || hcisResponseStatus.RetCode == "317" {
		// validate if username already exist
		if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ? AND application_code = ?", newUserRequest.Username, appDetails.Application_code).Scan(&userDetailValidation).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if userDetailValidation.User_id != 0 {
			returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Already Exists", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// validate if phone number already exist
		if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE phone_no = ? AND application_code = ?", isPhoneNoFormatted.Data.Message, appDetails.Application_code).Scan(&userDetailValidation).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if userDetailValidation.User_id != 0 {
			returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Phone Number Already Exists", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", newUserRequest.Institution_code).Scan(&instiDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if instiDetails.Institution_id == 0 {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Code Not Foound", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// register the user
		if insertErr := database.DBConn.Raw("SELECT public.register_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, newUserRequest.First_name, newUserRequest.Middle_name, newUserRequest.Last_name, newUserRequest.Email, isPhoneNoFormatted.Data.Message, newUserRequest.Staff_id, instiDetails.Institution_id, hashTempPassword, true, "", isBdateFormatted.Data.Message, instiDetails.Institution_code, appDetails.Application_code, appDetails.Application_id).Scan(&remark).Error; insertErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, insertErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, remark)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	} else if !hcisResponseStatus.Data.IsSuccess {
		return c.JSON(hcisResponseStatus)
	} else {
		// register the user
		if insertErr := database.DBConn.Raw("SELECT public.register_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, hcisResponseDeatails.First_name, hcisResponseDeatails.Middle_name, hcisResponseDeatails.Last_name, hcisResponseDeatails.Email, hcisResponseDeatails.Phone_no, hcisResponseDeatails.Staff_id, hcisResponseDeatails.Institution_id, hashTempPassword, true, "", isBdateFormatted.Data.Message, hcisResponseDeatails.Institution_code, appDetails.Application_code, appDetails.Application_id).Scan(&remark).Error; insertErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, insertErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, remark)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	// get user details
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE (staff_id = ? OR username = ?) AND application_code = ?", newUserRequest.Staff_id, newUserRequest.Username, appDetails.Application_code).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "User Added Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// append the temp password to the user's details
	UserDetails.Password = tempPassword

	// marshal the response body
	UserDetailsByte, marshalErr := json.Marshal(UserDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successResp := middleware.ResponseData(UserDetails.Username, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Registered User", nil, UserDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)
}
