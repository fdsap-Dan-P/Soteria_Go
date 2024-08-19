package usermanagement

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/controllers/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserFromHCIS(c *fiber.Ctx) error {
	registrationRequest := request.RegisterUserRequest{}
	userStaffId := request.ValidateInHCIS{}
	userHCISInfo := response.StaffInfoResponse{}
	instiDetails := response.InstitutionDetails{}
	remark := response.DBFuncResponse{}
	userAccountDetails := response.UserAccountResponse{}
	userStatusDetails := response.UserStatusResponse{}

	currentDateTime := middleware.GetDateTime().Data.Message

	methodUsed := c.Method()
	endpoint := c.Path()
	funcName := "User Management"

	userActivity := "Register User"

	// // Extract JWT token from Authorization header
	// authHeader := c.Get("Authorization")
	// token := strings.TrimPrefix(authHeader, "Bearer")
	// tokenString := strings.TrimSpace(token)

	// if strings.TrimSpace(authHeader) == "" || tokenString == "" {
	// 	returnMessage := middleware.ResponseData("", funcName, "111", methodUsed, endpoint, []byte(""), []byte(""), "", nil)
	// 	if !returnMessage.Data.IsSuccess {
	// 		return c.JSON(returnMessage)
	// 	}
	// }

	// // Validate JWT token
	// claims := middleware.ParseToken(tokenString, funcName, methodUsed, endpoint)
	// if !claims.Data.IsSuccess {
	// 	return c.JSON(claims)
	// }

	// // get username
	// username := claims.Data.Message
	username := registrationRequest.Username

	// get request body
	if parsErr := c.BodyParser(&registrationRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Initial Request Body Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshhall the initial request body
	registrationRequestByte, marshallErr := json.Marshal(registrationRequest)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, registrationRequestByte, []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if staff id was filled
	if strings.TrimSpace(registrationRequest.Staff_id) == "" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, registrationRequestByte, []byte(""), "Employee ID Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validated the staff id format
	isSatffIdValid := validations.ValidateStaffID(registrationRequest.Staff_id)
	if !isSatffIdValid {
		returnMessage := middleware.ResponseData(username, funcName, "112", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fmt.Errorf(userStaffId.StaffID))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if staff id is already registered
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE staff_id = ?", registrationRequest.Staff_id).Scan(&userAccountDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userAccountDetails.User_id != 0 {
		if fetchErr := database.DBConn.Raw("SELECT * FROM user_status WHERE status_id = ?", userAccountDetails.Status_id).Scan(&userStatusDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if userStatusDetails.Status_id != 0 {
			returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, registrationRequestByte, []byte(""), "User Status Not Found", fmt.Errorf(userStaffId.StaffID))
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if userStatusDetails.Status_name != "Active" {
			errMessage := "Staff ID Already Registered | " + userStatusDetails.Status_name + " User Account"
			returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, registrationRequestByte, []byte(""), errMessage, fmt.Errorf(userStaffId.StaffID))
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		returnMessage := middleware.ResponseData(username, funcName, "403", methodUsed, endpoint, registrationRequestByte, []byte(""), "Staff ID Already Registered", fmt.Errorf(userStaffId.StaffID))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(registrationRequest.Username) == "" {
		registrationRequest.Username = registrationRequest.Staff_id
	}

	userStaffId.StaffID = registrationRequest.Staff_id
	userStaffIdByte, marshalErr := json.Marshal(userStaffId)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "311", methodUsed, endpoint, registrationRequestByte, []byte(""), "", marshalErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	// Basic auth credentials
	hcis_username := "fdsap_apis"
	hcis_password := "P@ssword123"
	hcis_authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(hcis_username+":"+hcis_password))

	// Create a new Resty client
	client := resty.New()

	// Send the request
	resp, respErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", hcis_authHeader).
		SetBody(userStaffIdByte).
		Post("https://ua-uat.cardmri.com:8555/HCISLink/WEBAPI/ExternalService/ViewStaffInfo")

	if respErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "317", methodUsed, endpoint, registrationRequestByte, []byte(""), "", respErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Unmarshal the response body into the struct
	if unmarshallErr := json.Unmarshal(resp.Body(), &userHCISInfo); unmarshallErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "310", methodUsed, endpoint, registrationRequestByte, []byte(""), "", unmarshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if len(userHCISInfo.StaffInfo) == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, registrationRequestByte, []byte(""), "Staff ID Not Found in HCIS", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Return the response struct as JSON
	for _, val := range userHCISInfo.StaffInfo {

		// data gathering for institution table
		if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_name = ?", val.Institution).Scan(&instiDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fetchErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// data gathering for institution table
		if instiDetails.Institution_id == 0 {
			if insErr := database.DBConn.Raw("INSERT INTO offices_mapping.institutions (institution_name) VALUES (?) RETURNING institution_id", val.Institution).Scan(&instiDetails).Error; insErr != nil {
				returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, registrationRequestByte, []byte(""), "", insErr)
				if !returnMessage.Data.IsSuccess {
					return c.JSON(returnMessage)
				}
			}
		}

		// insert data into user table
		if insertErr := database.DBConn.Raw("SELECT register_user(?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", registrationRequest.Username, val.FirstName, val.LastName, val.EmailAddress, val.MobilePhone, val.StaffID, 1, currentDateTime, instiDetails.Institution_id).Scan(&remark).Error; insertErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, registrationRequestByte, []byte(""), "", insertErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if remark.Remark != "Success" {
			returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fmt.Errorf(remark.Remark))
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	// get user_id
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE staff_id = ?", registrationRequest.Staff_id).Scan(&userAccountDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// generate user's temporary password
	temporaryPassword := middleware.PasswordGeneration()
	hashTemporaryPassword := hash.SHA256(temporaryPassword)

	// store user's temporary password
	if insertErr := database.DBConn.Raw("SELECT add_user_passwords(?, ?, ?, ?) AS remark", userAccountDetails.User_id, hashTemporaryPassword, true, currentDateTime).Scan(&remark).Error; insertErr != nil {
		if deletErr := database.DBConn.Raw("DELETE FROM user_accounts WHERE user_id = ?", userAccountDetails.User_id).Scan(&userAccountDetails).Error; deletErr != nil {
			returnMessage := middleware.ResponseData(username, funcName, "314", methodUsed, endpoint, registrationRequestByte, []byte(""), "", deletErr)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, registrationRequestByte, []byte(""), "", insertErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(username, funcName, "303", methodUsed, endpoint, registrationRequestByte, []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	userFullName := userAccountDetails.First_name + " " + userAccountDetails.Last_name
	mailBody := "Dear " + userFullName + ", \r\n\nHere is your New Cagabay Account credentials. \n\nUsername: " + userAccountDetails.Username + "\nor\nStaff ID: " + userAccountDetails.Staff_id + "\n\nTemporary Password: " + temporaryPassword + "\n\n\n\nYou can login here:\nhttps://bakawan-rbi.fortress-asya.com\n\n Thank you, \n\nData Platform Support Team"

	// send email
	mailResponse := middleware.SendMail(userFullName, userAccountDetails.Email, "Cagabay New User Account: Temporary Credentials", mailBody, userAccountDetails.Username, funcName, methodUsed, endpoint, registrationRequestByte, []byte(""))
	if !mailResponse.Data.IsSuccess {
		return c.JSON(mailResponse)
	}
	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(registrationRequest.Username, username, userActivity, string(registrationRequestByte), "Registered User", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(username, "User Log", "202", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Log Out", "", nil)

	return c.JSON(response.ResponseModel{
		RetCode: "203",
		Message: "Successfully Registered!",
	})
}
