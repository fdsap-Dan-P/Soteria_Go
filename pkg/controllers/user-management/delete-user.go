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

func DeleteUser(c *fiber.Ctx) error {
	userIdentifier := request.LoginCredentialsRequest{}
	userDetails := response.UserDetails{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Delete User"

	// Extraxt the headers
	apiKey := c.Get("X-API-Key")
	authHeader := c.Get("Authorization")

	validationStatus, validationDetails := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !validationStatus.Data.IsSuccess {
		return c.JSON(validationStatus)
	}

	// parse the request body
	if parsErr := c.BodyParser(&userIdentifier); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", validationDetails.App_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	userIdentifierByte, marshalErr := json.Marshal(userIdentifier)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(userIdentifier.User_identity, "", validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userIdentifier.User_identity) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "401", methodUsed, endpoint, userIdentifierByte, []byte(""), "User Identity Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userIdentifier.Institution_code) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "401", methodUsed, endpoint, userIdentifierByte, []byte(""), "User Institution Code Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_details WHERE (username = ? OR staff_id = ? OR email = ? OR phone_no = ?) AND institution_code = ? AND application_code = ?", userIdentifier.User_identity, userIdentifier.User_identity, userIdentifier.User_identity, userIdentifier.User_identity, userIdentifier.Institution_code, validationDetails.App_code).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "302", methodUsed, endpoint, userIdentifierByte, []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "404", methodUsed, endpoint, userIdentifierByte, []byte(""), "User To Be Deleted Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if deletErr := database.DBConn.Raw("DELETE FROM user_accounts WHERE user_id = ?", userDetails.User_id).Scan(&userDetails).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "314", methodUsed, endpoint, userIdentifierByte, []byte(""), "", deletErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "210", methodUsed, endpoint, userIdentifierByte, []byte(""), "", nil, userDetails)

	return c.JSON(successMessage)
}
