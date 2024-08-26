package securitymanagement

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AppRegistration(c *fiber.Ctx) error {
	newAppRequest := request.ApplicationRequest{}
	appDetails := response.ApplicationDetails{}

	moduleName := "Security Management"
	funcName := "Application Registration"
	methodUsed := c.Method()
	endpoint := c.Path()

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
	if fetchErr := database.DBConn.Debug().Raw("SELECT * FROM public.applications WHERE application_name = ?", newAppRequest.App_name).Scan(&appDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "302", methodUsed, endpoint, newAppRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fmt.Println(appDetails.Application_id)

	if appDetails.Application_id != 0 {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "403", methodUsed, endpoint, newAppRequestByte, []byte(""), "Application Name Already Exists", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// generate the app code
	appCode, genErr := middleware.AppCodeGeneration(newAppRequest.App_name)
	if genErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "302", methodUsed, endpoint, newAppRequestByte, []byte(""), "", genErr, genErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// generate the api key
	apiKey := uuid.New().String()

	// insert the app details
	if insErr := database.DBConn.Raw("INSERT INTO public.applications(application_code, application_name, application_description, api_key) VALUES (?, ?, ?, ?) RETURNING *", appCode, newAppRequest.App_name, newAppRequest.App_desc, apiKey).Scan(&appDetails).Error; insErr != nil {
		returnMessage := middleware.ResponseData("", "", "", moduleName, funcName, "303", methodUsed, endpoint, newAppRequestByte, []byte(""), "", genErr, genErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// remove app id as return
	appDetails.Application_id = 0

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
