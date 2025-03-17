package userlogs

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

func Login(c *fiber.Ctx) error {
	credentialRequest := request.LoginCredentialsRequest{}
	userDetails := response.UserDetails{}
	userPasswordDetails := response.UserPasswordDetails{}
	userAppDetails := response.UserApplicationDetails{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Logs"
	funcName := "Log In"

	// Extraxt the api key
	apiKey := c.Get("X-API-Key")
	fmt.Println("API KEY: ", apiKey)

	// validate the api key
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !apiKeyValidatedStatus.Data.IsSuccess {
		return c.JSON(apiKeyValidatedStatus)
	}

	// parse the request body
	if parsErr := c.BodyParser(&credentialRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", appDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	credentialRequestByte, marshalErr := json.Marshal(credentialRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if staff_id has value
	if strings.TrimSpace(credentialRequest.User_identity) == "" {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Unique Identity Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if staff_id has value
	if strings.TrimSpace(credentialRequest.Password) == "" {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Password Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if user identity is valid
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.user_details WHERE staff_id = ? OR username = ? OR email = ? OR phone_no = ?", credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, credentialRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if user is linked the application
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_app_view WHERE (username = ? OR staff_id = ? OR email = ? OR phone_no = ?) AND application_code = ?", credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity, appDetails.Application_code).Scan(&userAppDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, credentialRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userAppDetails.Username) == "" { // user is not linked to the application
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "402", methodUsed, endpoint, credentialRequestByte, []byte(""), "Unauthorized Access", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if password is valid
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_passwords WHERE user_id = ? AND insti_code = ? AND app_code = ? ORDER BY created_at DESC LIMIT 1", userDetails.User_id, userDetails.Institution_code, appDetails.Application_code).Scan(&userPasswordDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, credentialRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userPasswordDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	hashPasswordRequest := hash.SHA256(credentialRequest.Password)
	if userPasswordDetails.User_id == 0 || userPasswordDetails.Password_hash != hashPasswordRequest {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "103", methodUsed, endpoint, credentialRequestByte, []byte(""), "Invalid Credential", nil, nil) // must no user details
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// append password details to user details
	userDetails.Requires_password_reset = userPasswordDetails.Requires_password_reset
	userDetails.Last_password_reset = userPasswordDetails.Last_password_reset

	fmt.Println("INSTI CODE: ", userDetails.Institution_code)

	// generate the jwt token
	token, tokenErr := middleware.GenerateToken(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint)
	fmt.Print("token: ", token)
	fmt.Println("tokenErr: ", tokenErr)
	if tokenErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "305", methodUsed, endpoint, credentialRequestByte, []byte(""), "", tokenErr, tokenErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// create or recrete the user token
	isTokenStored := middleware.StoringUserToken(token, userDetails.Username, userDetails.Staff_id, userDetails.Institution_code, appDetails.Application_code, moduleName, methodUsed, endpoint, credentialRequestByte)
	if !isTokenStored.Data.IsSuccess {
		return c.JSON(isTokenStored)
	}

	userDetails.Token = token

	// marshal the user details
	userDetailsByte, marshalErr := json.Marshal(userDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, credentialRequestByte, []byte(""), "Marshalling Response Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fmt.Println("userDetailsByte: ", string(userDetailsByte))

	returnMessage := middleware.ResponseData(credentialRequest.User_identity, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "201", methodUsed, endpoint, credentialRequestByte, userDetailsByte, "", nil, userDetails)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}

	return c.JSON(returnMessage)
}
