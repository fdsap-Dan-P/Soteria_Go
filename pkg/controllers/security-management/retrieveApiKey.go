package securitymanagement

import (
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/encryptDecrypt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func RetrievePlainApiKey(c *fiber.Ctx) error {
	appCode := c.Params("app-code")
	appDetails := response.ApplicationDetails{}
	userApiKeyDetails := response.ApplicationDetails{}

	moduleName := "Security Management"
	funcName := "Retrieve Plain API Key"
	methodUsed := c.Method()
	endpoint := c.Path()

	if strings.TrimSpace(appCode) == "" {
		returnMessage := middleware.ResponseData("", "", "", "", "", "401", c.Method(), c.Path(), []byte(""), []byte(""), "App Code Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if user is admin
	if appCode == "CU0003-1738988675" { // 1189289a-e743-46ee-87b0-53d2a08386b6
		// Extract the api key from Authorization header
		userApiKey := c.Get("X-API-Key")

		// validate if the api key is provided
		if strings.TrimSpace(userApiKey) == "" {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "User API Key Missing", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// get the secret key
		secretKey := os.Getenv("SECRET_KEY")
		if strings.TrimSpace(secretKey) == "" {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "Secret Key Not Found in Environment", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// encrypt the user api key
		userEncryptedApiKey, userEncryptErr := encryptDecrypt.EncryptWithSecretKey(userApiKey, secretKey)
		if userEncryptErr != nil {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "318", methodUsed, endpoint, []byte(""), []byte(""), "Encrypting User API Key Failed", userEncryptErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// get the user api key details
		if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.applications WHERE api_key = ?", userEncryptedApiKey).Scan(&userApiKeyDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// validate the user api key exists
		if userApiKeyDetails.Application_id == 0 {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "User API Key Not Found", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// validate if the user is an admin
		if userApiKeyDetails.Application_code != "CU0003-1738988675" {
			returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "402", methodUsed, endpoint, []byte(""), []byte(""), "Unauthorized Access", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE app_code = ?", appCode).Scan(&appDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData("", "", "", "", "", "302", c.Method(), c.Path(), []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if appDetails.Application_id == 0 {
		returnMessage := middleware.ResponseData("", "", "", "", "", "404", c.Method(), c.Path(), []byte(""), []byte(""), "Application Not Found", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get the secret key
	secretKey := os.Getenv("SECRET_KEY")
	if strings.TrimSpace(secretKey) == "" {
		returnMessage := middleware.ResponseData("", "", "", "", "", "404", c.Method(), c.Path(), []byte(""), []byte(""), "Secret Key Not Found in Environment", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// decrypt the api key
	decryptedApiKey, decryptErr := encryptDecrypt.DecryptWithSecretKey(appDetails.Api_key, secretKey)
	if decryptErr != nil {
		returnMessage := middleware.ResponseData("", "", "", "", funcName, "319", c.Method(), c.Path(), []byte(""), []byte(""), "Decrypting API Key Failed", decryptErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	appDetails.Api_key = decryptedApiKey
	return c.JSON(appDetails)
}
