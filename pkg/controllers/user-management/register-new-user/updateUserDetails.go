package registernewuser

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UpdateUserDetails(c *fiber.Ctx) error {
	userCategory := c.Params("user_category")
	userIdentity := c.Params("user_identity")
	newUserRequest := request.UserRegistrationRequest{}
	UserDetails := response.UserDetails{}
	updatedUserDetails := response.UserDetails{}
	UserDetailsChecker := response.UserDetails{}
	institutionDetails := response.InstitutionDetails{}
	remark := response.DBFuncResponse{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Register New User"

	// Extraxt the headers
	apiKey := c.Get("X-API-Key")
	authHeader := c.Get("Authorization")

	validationStatus, validationDetails := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !validationStatus.Data.IsSuccess {
		return c.JSON(validationStatus)
	}

	// check if to be updated exist
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE (staff_id = ? OR username = ? OR email = ? OR phone_no = ?) AND application_code = ?", userIdentity, userIdentity, userIdentity, userIdentity, validationDetails.App_code).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	userIdToBeUpdated := UserDetails.User_id
	if validationDetails.Insti_code != UserDetails.Institution_code {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "402", methodUsed, endpoint, []byte(""), []byte(""), "Unauthorized To Update This User", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// parse the request body
	if parsErr := c.BodyParser(&newUserRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	newUserRequestByte, marshalErr := json.Marshal(newUserRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userCategory == "staff" {
		// check if staff_id has value
		if strings.TrimSpace(newUserRequest.Staff_id) == "" {
			returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Missing", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if strings.TrimSpace(newUserRequest.Username) == "" {
			newUserRequest.Username = newUserRequest.Staff_id
		}

		// validate the staff id format
		isStaffIdValidated := validations.StaffIdValidation(newUserRequest.Staff_id, moduleName, methodUsed, endpoint)
		if !isStaffIdValidated {
			returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "112", methodUsed, endpoint, newUserRequestByte, []byte(""), "Invalid Employee ID", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// validate if username already exists
		if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE staff_id = ? AND user_id != ? AND application_code = ? AND institution_code = ?", newUserRequest.Staff_id, userIdToBeUpdated, validationDetails.App_code, validationDetails.Insti_code).Scan(&UserDetailsChecker).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if UserDetailsChecker.User_id != 0 {
			returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Already Exists", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	} else if userCategory == "non-staff" {
		if strings.TrimSpace(newUserRequest.Username) == "" {
			returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Missing", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	} else {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "100", methodUsed, endpoint, newUserRequestByte, []byte(""), "Validation Failed | Invalid User Category", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Institution_code) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Code Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Phone_no) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Phone Number Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Birthdate) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Brthd Date Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Email) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Email Address Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// format the phone number
	isPhoneNoFormatted := middleware.NormalizePhoneNumber(newUserRequest.Phone_no, newUserRequest.Username, validationDetails.Insti_code, validationDetails.App_code, funcName, methodUsed, endpoint)
	if !isPhoneNoFormatted.Data.IsSuccess {
		return c.JSON(isPhoneNoFormatted)
	}

	// format the birthdate
	isBdateFormatted := middleware.FormatingDate(newUserRequest.Birthdate, newUserRequest.Username, validationDetails.Insti_code, validationDetails.App_code, funcName, methodUsed, endpoint)
	if !isBdateFormatted.Data.IsSuccess {
		return c.JSON(isBdateFormatted)
	}

	// validate email address
	isEmailAddrValid := middleware.ValidateEmail(newUserRequest.Email)
	if !isEmailAddrValid {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "109", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if username already exists
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE username = ? AND user_id != ? AND application_code = ? AND institution_code = ?", newUserRequest.Username, userIdToBeUpdated, validationDetails.App_code, validationDetails.Insti_code).Scan(&UserDetailsChecker).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetailsChecker.User_id != 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if phone number already exists
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE phone_no = ? AND user_id != ? AND application_code = ? AND institution_code = ?", newUserRequest.Phone_no, userIdToBeUpdated, validationDetails.App_code, validationDetails.Insti_code).Scan(&UserDetailsChecker).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetailsChecker.User_id != 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Phone Number Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if email address already exists
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE email = ? AND user_id != ? AND application_code = ? AND institution_code = ?", newUserRequest.Email, userIdToBeUpdated, validationDetails.App_code, validationDetails.Insti_code).Scan(&UserDetailsChecker).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetailsChecker.User_id != 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Email Address Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", newUserRequest.Institution_code).Scan(&institutionDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if institutionDetails.Institution_id == 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// register the user
	if insertErr := database.DBConn.Raw("SELECT public.update_user_info(?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, newUserRequest.First_name, newUserRequest.Middle_name, newUserRequest.Last_name, newUserRequest.Email, isPhoneNoFormatted.Data.Message, newUserRequest.Staff_id, institutionDetails.Institution_id, userIdToBeUpdated).Scan(&remark).Error; insertErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "304", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "304", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fmt.Errorf("%s", remark.Remark), nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user details
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ? AND application_code = ? AND institution_code = ?", newUserRequest.Username, validationDetails.App_code, newUserRequest.Institution_code).Scan(&updatedUserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the response body
	UserDetailsByte, marshalErr := json.Marshal(updatedUserDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successResp := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "204", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Updated User Details", nil, updatedUserDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	return c.JSON(successResp)
}
