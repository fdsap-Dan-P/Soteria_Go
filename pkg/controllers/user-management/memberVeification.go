package usermanagement

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	dmUserDetails := response.MemberResponse{}
	dmUserSavings := response.MemberResponse{}
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

	if strings.TrimSpace(userRequest.Phone_no) == "" && strings.TrimSpace(userRequest.First_name) == "" && strings.TrimSpace(userRequest.Last_name) == "" && strings.TrimSpace(userRequest.Birthdate) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "401", methodUsed, endpoint, userRequestByte, []byte(""), "Phone Number Input Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// validate the phone number format
	isPhoneNumberValidated := middleware.NormalizePhoneNumber(userRequest.Phone_no, fullName, "", validationDetails.Application_code, funcName, methodUsed, endpoint)
	if !isPhoneNumberValidated.Data.IsSuccess {
		return c.JSON(isPhoneNumberValidated)
	}
	userRequest.Phone_no = isPhoneNumberValidated.Data.Message

	// format the birthdate
	isBirthDateValid := middleware.FormatingDate(userRequest.Birthdate, fullName, "", validationDetails.Application_code, moduleName, funcName, methodUsed, endpoint)
	if !isBirthDateValid.Data.IsSuccess {
		return c.JSON(isBirthDateValid)
	}
	userRequest.Birthdate = isBirthDateValid.Data.Message

	// validate user to data mart
	// get API key
	apiKey := os.Getenv("DATA_MART_API_KEY")
	if strings.TrimSpace(apiKey) == "" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "404", methodUsed, endpoint, userRequestByte, []byte(""), "API KEY NOT FOUND IN ENVIRONMENT", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	dmMemberVerificationUrl := os.Getenv("DATA_MART_HOST") + "/users/verify-member"

	dmMemberVerificationtr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	dmMemberVerificationClient := &http.Client{Transport: dmMemberVerificationtr}

	// marshal the request body for data mart
	dmMemberVerificationRequestByte, marshalErr := json.Marshal(userRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, userRequestByte, []byte(""), "Marshalling Request Body Failed", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// Create the HTTP request and set headers
	dmMemberVerificationReq, dmMemberVerificationReqErr := http.NewRequest("POST", dmMemberVerificationUrl, bytes.NewBuffer(dmMemberVerificationRequestByte))
	if dmMemberVerificationReqErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "318", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberVerificationReqErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	dmMemberVerificationReq.Header.Set("Content-Type", "application/json")
	dmMemberVerificationReq.Header.Set("X-API-Key", apiKey)

	fmt.Println(dmMemberVerificationReq)

	// Send the request
	dmMemberVerificationResp, dmMemberVerificationRespErr := dmMemberVerificationClient.Do(dmMemberVerificationReq)
	if dmMemberVerificationRespErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "317", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberVerificationRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	defer dmMemberVerificationResp.Body.Close()

	if dmMemberVerificationResp.Status != "200 OK" {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "405", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberVerificationRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	dmMemberVerificationBody, dmMemberVerificationErr := ioutil.ReadAll(dmMemberVerificationResp.Body)
	if dmMemberVerificationErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "317", methodUsed, endpoint, userRequestByte, []byte(""), "Reading Data Mart Response Failed", dmMemberVerificationErr, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// soteriaResponse_v2 := map[string]string{}
	if err := json.Unmarshal(dmMemberVerificationBody, &dmUserDetails); err != nil {
		returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "310", methodUsed, endpoint, userRequestByte, []byte(""), "", err, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	userVerification.Phone_no = userRequest.Phone_no
	userVerification.First_name = userRequest.First_name
	userVerification.Last_name = userRequest.Last_name
	userVerification.Birthdate = userRequest.Birthdate

	if strings.TrimSpace(dmUserDetails.Data.Details.Cid) == "" {
		userVerification.Is_member = false
	} else {
		userVerification.Is_member = true
		memberDetails["member_details"] = dmUserDetails.Data.Details

		// get member's institution details
		if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", dmUserDetails.Data.Details.Insti_code).Scan(&instiDetails).Error; fetchErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "310", methodUsed, endpoint, userRequestByte, []byte(""), "", fetchErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		userVerification.Institution_code = instiDetails.Institution_code
		userVerification.Institution_name = instiDetails.Institution_name

		// make request to data mart to  get user's saving account
		dmMemberSavingRequestBody := request.MemberVerificationRequest{
			Cid:              dmUserDetails.Data.Details.Cid,
			Institution_code: dmUserDetails.Data.Details.Insti_code,
		}
		// get member's saving account details
		dmMemberSavingUrl := os.Getenv("DATA_MART_HOST") + "/users/saving-account"

		dmMemberSavingtr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		dmMemberSavingClient := &http.Client{Transport: dmMemberSavingtr}

		// marshal the request body for data mart
		dmMemberSavingRequestByte, marshalErr := json.Marshal(dmMemberSavingRequestBody)
		if marshalErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "311", methodUsed, endpoint, userRequestByte, []byte(""), "Marshalling Request Body Failed", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// Create the HTTP request and set headers
		dmMemberSavingReq, dmMemberSavingReqErr := http.NewRequest("POST", dmMemberSavingUrl, bytes.NewBuffer(dmMemberSavingRequestByte))
		if dmMemberVerificationReqErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "318", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberSavingReqErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		dmMemberSavingReq.Header.Set("Content-Type", "application/json")
		dmMemberSavingReq.Header.Set("X-API-Key", apiKey)

		fmt.Println(dmMemberSavingReq)

		// Send the request
		dmMemberSavingResp, dmMemberSavingRespErr := dmMemberSavingClient.Do(dmMemberSavingReq)
		if dmMemberVerificationRespErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "317", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberSavingRespErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}
		defer dmMemberSavingResp.Body.Close()

		if dmMemberVerificationResp.Status != "200 OK" {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "405", methodUsed, endpoint, userRequestByte, []byte(""), "", dmMemberSavingRespErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		dmMemberSavingBody, dmMemberSavingErr := ioutil.ReadAll(dmMemberVerificationResp.Body)
		if dmMemberVerificationErr != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "317", methodUsed, endpoint, userRequestByte, []byte(""), "Reading Data Mart Response Failed", dmMemberSavingErr, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		// soteriaResponse_v2 := map[string]string{}
		if err := json.Unmarshal(dmMemberSavingBody, &dmUserSavings); err != nil {
			returnMessage := middleware.ResponseData(fullName, "", validationDetails.Application_code, moduleName, funcName, "310", methodUsed, endpoint, userRequestByte, []byte(""), "", err, nil)
			if !returnMessage.Data.IsSuccess {
				return c.JSON(returnMessage)
			}
		}

		memberDetails["saving_details"] = dmUserSavings.Data.Details
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
