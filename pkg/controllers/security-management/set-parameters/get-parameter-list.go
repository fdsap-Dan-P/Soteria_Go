package setparameters

import (
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func ParameterList(c *fiber.Ctx) error {
	moduleName := "Security Management"
	funcName := "Set JWT Parameter"
	methodUsed := c.Method()
	endpoint := c.Path()
	configParam := []response.ConfigDetails{}

	// extract headers
	authHeader := c.Get("Authorization")
	apiKey := c.Get("X-API-Key")

	// validate the header
	headerValidationStatus, headerValidationResponse := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !headerValidationStatus.Data.IsSuccess {
		return c.JSON(headerValidationStatus)
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config_params WHERE config_insti_code = ? AND config_app_code = ?", headerValidationResponse.Insti_code, headerValidationResponse.App_code).Scan(&configParam).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if len(configParam) == 0 {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", nil, configParam)

	return c.JSON(successMessage)
}
