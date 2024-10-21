package usermanagement

import (
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteUser(c *fiber.Ctx) error {
	userIdentity := c.Params("user_identity")
	userDetails := response.UserDetails{}
	remark := response.DBFuncResponse{}

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

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_details WHERE username = ? OR staff_id = ?", userIdentity, userIdentity).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User To Be Deleted Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if deletErr := database.DBConn.Raw("DELETE FROM user_accounts WHERE username = ? OR staff_id = ? AS remark", userIdentity, userIdentity).Scan(&remark).Error; deletErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "314", methodUsed, endpoint, []byte(""), []byte(""), "", deletErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successMessage := middleware.ResponseData(validationDetails.Username, validationDetails.App_code, validationDetails.Insti_code, moduleName, funcName, "210", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)

	return c.JSON(successMessage)
}
