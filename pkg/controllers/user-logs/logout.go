package userlogs

import (
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
)

func LogOut(c *fiber.Ctx) error {
	username := c.Params("username")
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

	if fetchErr := database.DBConn.Raw("SELECT * FROM public.user_details WHERE username = ?", username).Scan(&userDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, "", appDetails.Application_code, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userDetails.User_id == 0 {
		returnMessage := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if deleteErr := database.DBConn.Exec("DELETE FROM logs.user_tokens WHERE username = ? OR staff_id = ?", username, username).Error; deleteErr != nil {
		returnMessage := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "314", methodUsed, endpoint, []byte(""), []byte(""), "", deleteErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	successResp := middleware.ResponseData(username, userDetails.Institution_code, appDetails.Application_code, moduleName, funcName, "202", methodUsed, endpoint, []byte(""), []byte(""), "User Logged Out Successfully", nil, userDetails)
	return c.JSON(successResp)
}
