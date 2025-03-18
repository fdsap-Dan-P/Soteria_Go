package userlogs

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

func LogOut(c *fiber.Ctx) error {
	credentialRequest := request.LoginCredentialsRequest{}
	userToken := response.UserTokenDetails{}

	moduleName := "User Logs"
	funcName := "Log Out"
	methodUsed := c.Method()
	endpoint := c.Path()

	userDetails := response.UserDetails{}

	// Extraxt the api key
	apiKey := c.Get("X-API-Key")

	// validate the api key
	apiKeyValidatedStatus, appDetails := validations.APIKeyValidation(apiKey, "", "", "", funcName, methodUsed, endpoint, []byte(""))
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

	// check if institution code has value
	if strings.TrimSpace(credentialRequest.Institution_code) == "" {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, "", appDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Institution Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE (staff_id = ? OR username = ? OR email = ? OR phone_no = ?) AND institution_code = ? AND application_code = ?", credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.User_identity, credentialRequest.Institution_code, appDetails.Application_code).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, credentialRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, credentialRequestByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(credentialRequest.User_identity, credentialRequest.Institution_code, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, credentialRequestByte, []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if deleteErr := database.DBConn.Raw("DELETE FROM logs.user_tokens WHERE username = ? OR staff_id = ?", userDetails.Username, userDetails.Staff_id).Scan(&userToken).Error; deleteErr != nil {
		returnMessage := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "314", methodUsed, endpoint, credentialRequestByte, []byte(""), "", deleteErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successResp := middleware.ResponseData(userDetails.Username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "202", methodUsed, endpoint, credentialRequestByte, []byte(""), "User Logged Out Successfully", nil, userDetails)
	return c.JSON(successResp)
}
