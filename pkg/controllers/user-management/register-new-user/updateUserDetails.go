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
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE staff_id = ? OR username = ? OR email = ? OR phone_no = ?", userIdentity, userIdentity, userIdentity, userIdentity).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if UserDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	userIdToBeUpdated := UserDetails.User_id

	// parse the request body
	if parsErr := c.BodyParser(&newUserRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", validationDetails.App_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	newUserRequestByte, marshalErr := json.Marshal(newUserRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if staff_id has value
	if strings.TrimSpace(newUserRequest.Staff_id) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Username) == "" {
		newUserRequest.Username = newUserRequest.Staff_id
	}

	if strings.TrimSpace(newUserRequest.Institution_code) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Code Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate the staff id format
	isStaffIdValidated := validations.StaffIdValidation(newUserRequest.Staff_id, moduleName, methodUsed, endpoint)
	if !isStaffIdValidated {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "112", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if staff id already exists
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE (username = ? OR staff_id = ?) AND user_id != ?", newUserRequest.Username, newUserRequest.Staff_id, userIdToBeUpdated).Scan(&UserDetailsChecker).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fmt.Println("UserDetails: ", UserDetails.User_id)

	if strings.TrimSpace(UserDetailsChecker.Staff_id) != "" {
		UserDetailsChecker.User_id = 0 // data privacy

		return c.JSON(response.ResponseModel{
			RetCode: "403",
			Message: "Validation Failed",
			Data: response.DataModel{
				Message:   "Username or Employee ID Already Exists",
				IsSuccess: false,
				Error:     nil,
				Details:   UserDetailsChecker,
			},
		})

		// send email to the user that statte
		// tried being added via <app name>
		// pleaase used the current credentials to login
		// or forget the password if you dont remember it
	}

	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", newUserRequest.Institution_code).Scan(&institutionDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if institutionDetails.Institution_id == 0 {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "Institution Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// register the user
	if insertErr := database.DBConn.Raw("SELECT public.update_user_info(?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, newUserRequest.First_name, newUserRequest.Middle_name, newUserRequest.Last_name, newUserRequest.Email, newUserRequest.Phone_no, newUserRequest.Staff_id, institutionDetails.Institution_id, userIdToBeUpdated).Scan(&remark).Error; insertErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "304", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, insertErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "304", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, remark)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get user details
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE staff_id = ? OR username = ?", newUserRequest.Staff_id, newUserRequest.Username).Scan(&updatedUserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", validationDetails.App_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
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

	successResp := middleware.ResponseData(UserDetails.Username, newUserRequest.Institution_code, validationDetails.App_code, moduleName, funcName, "204", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Updated User Details", nil, updatedUserDetails)
	if !successResp.Data.IsSuccess {
		return c.JSON(successResp)
	}

	fmt.Println(successResp)

	return c.JSON(successResp)
}
