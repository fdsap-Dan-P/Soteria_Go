package memberVerification

import (
	"encoding/json"
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/middleware/validations"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MemberVerification(c *fiber.Ctx) error {
	userRequest := request.MemberVerificationRequest{}
	// userDetails := response.UserDetails{}
	instiDetails := response.InstitutionDetails{}

	userVerification := response.MemberVerificationResponse{}
	// dmUserDetails := response.MemberResponse{}
	// dmUserSavings := response.MemberResponse{}
	memberDetails := make(map[string]interface{})

	methodUsed := c.Method()
	endpoint := c.Path()
	moduleName := "User Management"
	funcName := "Member Verification"

	// Extraxt the headers
	requesterApiKey := c.Get("X-API-Key")

	validationStatus, validationDetails := validations.APIKeyValidation(requesterApiKey, "", "", "", moduleName, methodUsed, endpoint, []byte(""))
	if !validationStatus.Data.IsSuccess {
		return c.JSON(validationStatus)
	}

	// get the request body
	if parsErr := c.BodyParser(&userRequest); parsErr != nil {
		returnMessage := middleware.ResponseData("", "", validationDetails.Application_code, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Request Body Failed", parsErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	fullName := "Full Name: " + userRequest.First_name + " " + userRequest.Last_name

	// marshal the request body
	userRequestByte, marshalErr := json.Marshal(userRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, userRequestByte, []byte(""), "Marshalling Request Body Failed", marshalErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.Phone_no) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "Phone Number Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(userRequest.Phone_no) != "" {
		// validate the phone number format
		isPhoneNumberValidated := middleware.NormalizePhoneNumber(userRequest.Phone_no, fullName, "", validationDetails.Application_code, funcName, methodUsed, endpoint)
		if !isPhoneNumberValidated.Data.IsSuccess {
			return c.JSON(isPhoneNumberValidated)
		}
		userRequest.Phone_no = isPhoneNumberValidated.Data.Message
	}

	// format the birthdate
	isBirthDateValid := middleware.FormatingDate(userRequest.Birthdate, fullName, "", validationDetails.Application_code, moduleName, methodUsed, endpoint)
	if !isBirthDateValid.Data.IsSuccess {
		return c.JSON(isBirthDateValid)
	}
	userRequest.Birthdate = isBirthDateValid.Data.Message

	// get API key
	apiKey := os.Getenv("DATA_MART_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "API KEY NOT FOUND IN ENVIRONMENT", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	isDMMemberStatus, isDMMemberDetails := VerifyMemberFromDataMart(fullName, validationDetails.Application_code, moduleName, methodUsed, endpoint, apiKey, userRequest)
	if !isDMMemberStatus.Data.IsSuccess {
		return c.JSON(isDMMemberStatus)
	}

	userVerification.Phone_no = userRequest.Phone_no
	userVerification.First_name = userRequest.First_name
	userVerification.Last_name = userRequest.Last_name
	userVerification.Birthdate = userRequest.Birthdate
	userVerification.No_phone_number_user = isDMMemberDetails.Data.No_phone_number_user

	if strings.TrimSpace(isDMMemberDetails.Data.Details.Cid) == "" {
		userVerification.Is_member = false
	} else {
		userVerification.Is_member = true
		memberDetails["member_details"] = isDMMemberDetails.Data.Details

		// get member's institution details
		if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", isDMMemberDetails.Data.Details.Insti_code).Scan(&instiDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "310", methodUsed, endpoint, userRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		userVerification.Institution_code = instiDetails.Institution_code
		userVerification.Institution_name = instiDetails.Institution_name

		isSavingAccGotStatus, isSavingAccGotDetails := GetMemberSavingAcountDetailsFromDataMart(isDMMemberDetails.Data.Details.Cid, isDMMemberDetails.Data.Details.Insti_code, fullName, validationDetails.Application_code, moduleName, methodUsed, endpoint, apiKey)
		if !isSavingAccGotStatus.Data.IsSuccess {
			return c.JSON(isSavingAccGotStatus)
		}

		memberDetails["saving_details"] = isSavingAccGotDetails.Data.Details
		userVerification.Member_details = memberDetails
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   "Successfully Verified",
			IsSuccess: true,
			Error:     nil,
			Details:   userVerification,
		},
	})
}
