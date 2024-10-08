package securitymanagement

import (
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"

	"github.com/gofiber/fiber/v2"
)

func ThirdPartyHeaderValidation(c *fiber.Ctx) error {
	moduleName := "Security Management"
	funcName := "Token Validation"
	methodUsed := c.Method()
	endpoint := c.Path()

	// Extract JWT token from Authorization header
	authHeader := c.Get("Authorization")
	apiKey := c.Get("X-API-Key")

	// validate the header
	headerValidationStatus, headerValidationResponse := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !headerValidationStatus.Data.IsSuccess {
		return c.JSON(headerValidationStatus)
	}

	// marshal the response
	headerValidationResponseByte, marshalErr := json.Marshal(headerValidationResponse)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Response Failed", marshalErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "215", methodUsed, endpoint, []byte(""), headerValidationResponseByte, "", nil, headerValidationResponse)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}

	return c.JSON(returnMessage)
}
