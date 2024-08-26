package setparameters

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

func SetParams(c *fiber.Ctx) error {
	configCode := c.Params("config_code")
	paramRequest := request.ParameterRequest{}
	paramDetaills := response.ConfigDetails{}
	paramParameter := response.ConfigDetails{}

	moduleName := "Security Management"
	funcName := "Set JWT Parameter"
	methodUsed := c.Method()
	endpoint := c.Path()

	// extract headers
	authHeader := c.Get("Authorization")
	apiKey := c.Get("X-API-Key")

	// validate the header
	headerValidationStatus, headerValidationResponse := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !headerValidationStatus.Data.IsSuccess {
		return c.JSON(headerValidationStatus)
	}

	// parse the request body
	if parsErr := c.BodyParser(&paramRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, parsErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	paramRequestByte, marshalErr := json.Marshal(paramRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if request body is not empty
	if strings.TrimSpace(paramRequest.No_of_minutes) == "" {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "401", methodUsed, endpoint, paramRequestByte, []byte(""), "Parameter Value Input Missing", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(paramRequest.Application_code) == "" {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "401", methodUsed, endpoint, paramRequestByte, []byte(""), "Application Code Input Missing", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(paramRequest.Institution_code) == "" {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "401", methodUsed, endpoint, paramRequestByte, []byte(""), "Institution Code Input Missing", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate if autholrized to set the parameter
	if headerValidationResponse.Insti_code != paramRequest.Institution_code || headerValidationResponse.App_code != paramRequest.Application_code {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "402", methodUsed, endpoint, paramRequestByte, []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// get jwt details
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE insti_code = ?", configCode).Scan(&paramDetaills).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, paramRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if paramDetaills.Config_id == 0 {
		errMessage := "Parameter " + configCode + " Not Found"
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "404", methodUsed, endpoint, paramRequestByte, []byte(""), errMessage, nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// set parameter
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.insti_app_config WHERE config_id = ? AND insti_code = ? AND app_code = ?", paramDetaills.Config_id, paramRequest.Institution_code, paramRequest.Application_code).Scan(&paramParameter).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, paramRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	retCode := ""
	if paramParameter.Config_id == 0 {
		if insErr := database.DBConn.Raw("INSERT INTO parameters.insti_app_config (config_id, config_value, insti_code, app_code) VALUES (?, ?, ?)", paramDetaills.Config_id, paramRequest.No_of_minutes, paramRequest.Institution_code, paramRequest.Application_code).Scan(&paramParameter).Error; insErr != nil {
			returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "303", methodUsed, endpoint, paramRequestByte, []byte(""), "", insErr, insErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		retCode = "203"
	} else {
		currentDateTime := middleware.GetDateTime().Data.Message
		if updatErr := database.DBConn.Raw("UPDATE parameters.insti_app_config SET config_value = ?, updated_at = ? WHERE config_id = ? AND insti_code = ? AND app_code = ?", paramRequest.No_of_minutes, currentDateTime, paramDetaills.Config_id, paramRequest.Institution_code, paramRequest.Application_code).Scan(&paramParameter).Error; updatErr != nil {
			returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "304", methodUsed, endpoint, paramRequestByte, []byte(""), "", updatErr, updatErr.Error())
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		retCode = "204"
	}

	// get the new value
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config_params WHERE config_id = ? AND insti_code = ? AND app_code = ?", paramDetaills.Config_id, paramRequest.Institution_code, paramRequest.Application_code).Scan(&paramDetaills).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "302", methodUsed, endpoint, paramRequestByte, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if paramDetaills.Config_id == 0 {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "404", methodUsed, endpoint, paramRequestByte, []byte(""), "New Value Not Found", nil, paramDetaills)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the response
	paramDetaillsByte, marshalErr := json.Marshal(paramDetaills)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, "311", methodUsed, endpoint, paramRequestByte, paramDetaillsByte, "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	returnMessage := middleware.ResponseData(headerValidationResponse.Username, headerValidationResponse.Insti_code, headerValidationResponse.App_code, moduleName, funcName, retCode, methodUsed, endpoint, paramRequestByte, paramDetaillsByte, "", nil, paramDetaills)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}

	return c.JSON(returnMessage)
}
