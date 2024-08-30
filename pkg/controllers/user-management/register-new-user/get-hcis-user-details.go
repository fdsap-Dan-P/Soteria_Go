package registernewuser

import (
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func HCISUserDetailsProvider(c *fiber.Ctx) error {
	inquiryRequest := request.UserRegistrationRequest{}

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "HCIS User Details Provider"

	// Extraxt the headers
	apiKey := c.Get("X-API-Key")
	authHeader := c.Get("Authorization")
	// validate the api key

	validationStatus, validationDetails := validations.HeaderValidation(authHeader, apiKey, moduleName, funcName, methodUsed, endpoint)
	if !validationStatus.Data.IsSuccess {
		return c.JSON(validationStatus)
	}

	// get the request body
	if parsErr := c.BodyParser(&inquiryRequest); parsErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, inquiryRequest)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal the request body
	inquiryRequestByte, marshalErr := json.Marshal(inquiryRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, inquiryRequestByte, []byte(""), "Marshalling Request Body Failed", marshalErr, inquiryRequest)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(inquiryRequest.Staff_id) == "" {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "401", methodUsed, endpoint, inquiryRequestByte, []byte(""), "Staff Id Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	HcisInquiryStatus, HcisInquiryDetails := HcisInquiry(inquiryRequest.Staff_id, validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, methodUsed, endpoint, inquiryRequestByte)
	if !HcisInquiryStatus.Data.IsSuccess {
		return c.JSON(HcisInquiryStatus)
	}

	// marshal the response
	HcisInquiryDetailsByte, marshalErr := json.Marshal(HcisInquiryDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "311", methodUsed, endpoint, inquiryRequestByte, []byte(""), "Marshalling Request Body Failed", marshalErr, HcisInquiryDetails)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	returnMessage := middleware.ResponseData(validationDetails.Username, validationDetails.Insti_code, validationDetails.App_code, moduleName, funcName, "200", methodUsed, endpoint, inquiryRequestByte, HcisInquiryDetailsByte, "", nil, HcisInquiryDetails)
	if !returnMessage.Data.IsSuccess {
		return c.JSON(returnMessage)
	}
	return c.JSON(returnMessage)
}
