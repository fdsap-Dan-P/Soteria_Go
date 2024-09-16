package registernewuser

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	newUserRequest := request.UserRegistrationRequest{}
	UserDetails := response.UserDetails{}
	remark := response.DBFuncResponse{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Register New User"

	// // Extraxt the api key
	apiKey := c.Get("X-API-Key")

	// validate the api key
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", funcName, methodUsed, endpoint, []byte(""))
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

	// validate the staff id format
	isStaffIdValidated := validations.StaffIdValidation(newUserRequest.Staff_id, moduleName, methodUsed, endpoint)
	if !isStaffIdValidated {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "112", methodUsed, endpoint, newUserRequestByte, []byte(""), "Staff ID Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if staff id already exists
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE staff_id = ?", newUserRequest.Staff_id).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fmt.Println("UserDetails: ", UserDetails.User_id)

	if strings.TrimSpace(UserDetails.Staff_id) != "" {
		UserDetails.User_id = 0 // data privacy

		return c.JSON(response.ResponseModel{
			RetCode: "403",
			Message: "Validation Failed",
			Data: response.DataModel{
				Message:   "Username or Employee ID Already Exists",
				IsSuccess: false,
				Error:     nil,
				Details:   UserDetails,
			},
		})

		// send email to the user that statte
		// tried being added via <app name>
		// pleaase used the current credentials to login
		// or forget the password if you dont remember it
	}

	// get hcis details
	hcisResponseStatus, hcisResponseDeatails := HcisInquiry(newUserRequest.Staff_id, newUserRequest.Username, "", "", "User Registration", methodUsed, endpoint, newUserRequestByte)
	if !hcisResponseStatus.Data.IsSuccess {
		return c.JSON(hcisResponseStatus)
	}

	// generate user's temp password
	tempPassword := middleware.PasswordGeneration()
	hashTempPassword := hash.SHA256(tempPassword)

	currentDateTime := middleware.GetDateTime().Data.Message
	// register the user
	if insertErr := database.DBConn.Raw("SELECT public.register_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, hcisResponseDeatails.First_name, hcisResponseDeatails.Middle_name, hcisResponseDeatails.Last_name, hcisResponseDeatails.Email, hcisResponseDeatails.Phone_no, newUserRequest.Staff_id, hcisResponseDeatails.Institution_id, hashTempPassword, true, currentDateTime).Scan(&remark).Error; insertErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, hcisResponseDeatails.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insertErr, insertErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, hcisResponseDeatails.Institution_code, appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fmt.Errorf(remark.Remark), remark)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if staff id already exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE staff_id = ? OR username = ?", newUserRequest.Staff_id, newUserRequest.Username).Scan(&UserDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the response body
	UserDetailsByte, marshalErr := json.Marshal(UserDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Staff_id, hcisResponseDeatails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// send confirmation that user was successfully registered with ther credentials
	userFullName := UserDetails.First_name + " " + UserDetails.Last_name
	mailBody := "Dear " + userFullName + ", \r\n\nHere is your New " + appDetails.Application_name + " Account credentials. \n\nUsername: " + UserDetails.Username + "\nor\n" + "Employee ID: " + UserDetails.Staff_id + "\n\n" + "Password: " + tempPassword + "\n\n\n\nYou can login here:\nhttps://bakawan-rbi.fortress-asya.com \n or via Cagabay Mobile App\n\n Thank you, \n\n" + appDetails.Application_name + " Support Team"
	sendEmailErr := middleware.SendMail(userFullName, UserDetails.Email, "New "+appDetails.Application_name+" Credentials", mailBody, UserDetails.Username, hcisResponseDeatails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint, newUserRequestByte, UserDetailsByte)
	if !sendEmailErr.Data.IsSuccess {
		return c.JSON(sendEmailErr)
	}

	returnMessage := middleware.ResponseData(UserDetails.Username, hcisResponseDeatails.Institution_code, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Registered User", nil, UserDetails)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}

	fmt.Println(returnMessage)

	return c.JSON(returnMessage)
}
