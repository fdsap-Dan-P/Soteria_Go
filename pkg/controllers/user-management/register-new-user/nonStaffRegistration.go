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
	userAppDetails := []response.UserApplicationDetails{}
	userAppResp := response.UserAppResponse{}
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

	if strings.TrimSpace(newUserRequest.Phone_no) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Phone Number Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(newUserRequest.Phone_no) == "" {
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, newUserRequestByte, []byte(""), "Birth Date Input Missing", nil, nil)
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
		returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "109", methodUsed, endpoint, newUserRequestByte, []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if username already exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", newUserRequest.Username).Scan(&userDetailValidation).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if phone number already exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE phone_no = ?", isPhoneNoFormatted.Data.Message).Scan(&userDetailValidation).Error; fetchErr != nil {
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

	// validate if email address already exist
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE email = ?", newUserRequest.Email).Scan(&userDetailValidation).Error; fetchErr != nil {
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
	if insertErr := database.DBConn.Raw("SELECT public.register_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) AS remark", newUserRequest.Username, newUserRequest.First_name, newUserRequest.Middle_name, newUserRequest.Last_name, newUserRequest.Email, isPhoneNoFormatted.Data.Message, "", instiDetails.Institution_id, hashTempPassword, true, "", isBdateFormatted.Data.Message).Scan(&remark).Error; insertErr != nil {
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
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", newUserRequest.Username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		fmt.Println("TRACE 1")
		fmt.Println("userDetails.User_id: ", userDetails.User_id)
		returnMessage := middleware.ResponseData(newUserRequest.Username, userInstiCode, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, newUserRequestByte, []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fmt.Println("APP ID: ", appDetails.Application_id)

	if userDetailValidation.User_id != 0 {
		fmt.Println("TRACE 2")
		fmt.Println("userDetails.User_id: ", userDetails.User_id)
		// validate where application did the user is linked
		if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_app_view WHERE username = ?", newUserRequest.Username).Scan(&userAppDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, newUserRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		if len(userAppDetails) == 0 { // if the user is not linked to any application
			// link the user to this application
			if insErr := database.DBConn.Raw("INSERT INTO public.user_applications (user_id, application_id) VALUES (?, ?)", userDetails.User_id, appDetails.Application_id).Scan(&userAppResp).Error; insErr != nil {
				returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insErr, nil)
				if !returnMessage.Data.IsSuccess {
					return c.JSON(returnMessage)
				}
			}

			userDetails.Password = "Use Current Password in Linked Application"

			// marshal the response body
			UserDetailsByte, marshalErr := json.Marshal(userDetails)
			if marshalErr != nil {
				returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
				if !returnMessage.Data.IsSuccess {
					return c.JSON(returnMessage)
				}
			}

			successResp := middleware.ResponseData(userDetails.Username, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Registered User", nil, userDetails)
			return c.JSON(successResp)

		} else { // check if the user is already linked to any application
			fmt.Println("TRACE 3")
			fmt.Println("LEN: ", len(userAppDetails))
			isUserLinked := false
			for _, userLinkedApp := range userAppDetails { // check if the user is already linked to this application
				if userLinkedApp.Application_code == appDetails.Application_code {
					isUserLinked = true
				}
			}

			if !isUserLinked { // if the user is not linked to this application
				fmt.Println("TRACE 4")
				// link the user to this application
				if insErr := database.DBConn.Raw("INSERT INTO public.user_applications (user_id, application_id) VALUES (?, ?)", userDetails.User_id, appDetails.Application_id).Scan(&userAppResp).Error; insErr != nil {
					returnMessage := middleware.ResponseData(newUserRequest.Staff_id, "", appDetails.Application_code, moduleName, funcName, "303", methodUsed, endpoint, newUserRequestByte, []byte(""), "", insErr, nil)
					if !returnMessage.Data.IsSuccess {
						return c.JSON(returnMessage)
					}
				}

				userDetails.Password = "Use Current Password in Linked Application"

				// marshal the response body
				UserDetailsByte, marshalErr := json.Marshal(userDetails)
				if marshalErr != nil {
					returnMessage := middleware.ResponseData(newUserRequest.Staff_id, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, newUserRequestByte, []byte(""), "", marshalErr, marshalErr.Error())
					if !returnMessage.Data.IsSuccess {
						return c.JSON(returnMessage)
					}
				}

				successResp := middleware.ResponseData(userDetails.Username, instiDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "203", methodUsed, endpoint, newUserRequestByte, UserDetailsByte, "Successfully Registered User", nil, userDetails)
				return c.JSON(successResp)
			} else { // if the user is already linked to this application
				fmt.Println("TRACE 5")
				returnMessage := middleware.ResponseData(newUserRequest.Username, "", appDetails.Application_code, moduleName, funcName, "403", methodUsed, endpoint, newUserRequestByte, []byte(""), "Username Already Exists", nil, nil)
				if !returnMessage.Data.IsSuccess {
					return c.JSON(returnMessage)
				}
			}
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

	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("successResp", successResp)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - -")

	return c.JSON(successResp)
}
