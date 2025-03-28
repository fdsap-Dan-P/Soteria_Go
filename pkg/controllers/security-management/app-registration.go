package securitymanagement

import (
	"encoding/json"
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/encryptDecrypt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AppRegistration(c *fiber.Ctx) error {
	newAppRequest := request.ApplicationRequest{}
	userApiKeyDetails := response.ApplicationDetails{}
	appDetails := response.ApplicationDetails{}

	moduleName := "Security Management"
	funcName := "Application Registration"
	methodUsed := c.Method()
	endpoint := c.Path()

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
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE api_key = ?", userEncryptedApiKey).Scan(&userApiKeyDetails).Error; fetchErr != nil {
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

	if parsErr := c.BodyParser(&newAppRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	newAppRequestByte, marshalErr := json.Marshal(newAppRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if app name was provided
	if strings.TrimSpace(newAppRequest.App_name) == "" {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "401", methodUsed, endpoint, newAppRequestByte, []byte(""), "Application Name Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// check if app name already exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications WHERE application_name = ?", newAppRequest.App_name).Scan(&appDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "302", methodUsed, endpoint, newAppRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if appDetails.Application_id != 0 {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "403", methodUsed, endpoint, newAppRequestByte, []byte(""), "Application Name Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// generate the app code
	appCode, genErr := middleware.AppCodeGeneration(newAppRequest.App_name)
	if genErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "302", methodUsed, endpoint, newAppRequestByte, []byte(""), "", genErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// generate the api key
	apiKey := uuid.New().String()
	encryptedApiKey, encryptErr := encryptDecrypt.EncryptWithSecretKey(apiKey, secretKey)
	if encryptErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "318", methodUsed, endpoint, newAppRequestByte, []byte(""), "Encrypting API Key Failed", encryptErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// insert the app details
	if insErr := database.DBConn.Raw("INSERT INTO public.applications(application_code, application_name, application_description, api_key) VALUES (?, ?, ?, ?) RETURNING *", appCode, newAppRequest.App_name, newAppRequest.App_desc, encryptedApiKey).Scan(&appDetails).Error; insErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "303", methodUsed, endpoint, newAppRequestByte, []byte(""), "", genErr, genErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// remove app id as return
	appDetails.Application_id = 0
	appDetails.Api_key = apiKey // return the plain api key

	// marshal the response
	appDetailsByte, marshalErr := json.Marshal(appDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "311", methodUsed, endpoint, newAppRequestByte, []byte(""), "Marshalling Response Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "203", methodUsed, endpoint, newAppRequestByte, appDetailsByte, "Successfully Registered Application", nil, appDetails)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}

	return c.JSON(returnMessage)
}
