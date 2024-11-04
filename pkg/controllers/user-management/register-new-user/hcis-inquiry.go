package registernewuser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/go-resty/resty/v2"
)

func HcisInquiry(staffId, username, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) (response.ReturnModel, response.UserDetails) {
	userHCISDetails := response.UserDetails{}
	userHCISInfo := response.StaffInfoResponse{}
	instiDetails := response.InstitutionDetails{}

	funcName := "HCIS Inquiry"

	if strings.TrimSpace(staffId) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, reqBody, []byte(""), "Staff Id Missing", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Basic auth credentials
	hcis_username := "fdsap_apis"
	hcis_password := "P@ssword123"
	hcis_authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(hcis_username+":"+hcis_password))

	hcis_reqBody := map[string]string{"StaffID": staffId}
	hcis_reqBodyByte, marshallErr := json.Marshal(hcis_reqBody)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "Marshalling HCIS Request Failed", marshallErr, hcis_reqBody)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Create a new Resty client
	client := resty.New()

	// Send the request
	resp, respErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", hcis_authHeader).
		SetBody(hcis_reqBodyByte).
		Post("https://ua-uat.cardmri.com:8555/HCISLink/WEBAPI/ExternalService/ViewStaffInfo")

	// check the response from HCIS
	fmt.Println("- - - - - - - - - - HCIS Response - - - - - - - - - - -")
	fmt.Println("STATUS: ", resp.StatusCode())
	fmt.Println("SUCCESS \n", resp)
	fmt.Println("ERROR: ", respErr)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	if respErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "317", methodUsed, endpoint, reqBody, []byte(""), "Reading HCIS Response Failed", respErr, resp)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Check if the request was successful (status code 200)
	fmt.Println("- - - - - STATUS CODE - - - - - - -")
	fmt.Println("resp.StatusCode: ", resp.StatusCode())
	fmt.Println("http.StatusOK: ", http.StatusOK)
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == 201 {
			returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "User Not Found", respErr, resp)
			if !returnMessage.Data.IsSuccess {
				return returnMessage, userHCISDetails
			}
		}

		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "405", methodUsed, endpoint, reqBody, []byte(""), "Request Failed To HCIS", respErr, resp)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}

	}

	// Unmarshal the response body into the struct
	if unmarshallErr := json.Unmarshal(resp.Body(), &userHCISInfo); unmarshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "310", methodUsed, endpoint, reqBody, []byte(""), "", unmarshallErr, unmarshallErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if len(userHCISInfo.StaffInfo) == 0 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "Staff ID Not Found in HCIS", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Return the response struct as JSON
	for _, val := range userHCISInfo.StaffInfo {
		userHCISDetails.First_name = val.FirstName
		userHCISDetails.Last_name = val.LastName
		userHCISDetails.Middle_name = val.MiddleName
		userHCISDetails.Email = val.EmailAddress
		userHCISDetails.Staff_id = val.StaffID
		userHCISDetails.Birthdate = val.BirthDate
		userHCISDetails.Phone_no = val.MobilePhone
		userHCISDetails.Institution_name = val.Institution
	}

	// check if user institution is already in database
	if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_name = ?", userHCISDetails.Institution_name).First(&instiDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if instiDetails.Institution_id == 0 {
		userHCISDetails.Institution_code = GenerateInstitutionCode(userHCISDetails.Institution_name, username, instiCode, appCode, moduleName, methodUsed, endpoint)
		if insErr := database.DBConn.Raw("INSERT INTO public.institutions (institution_name, institution_code) VALUES (?, ?) RETURNING *", userHCISDetails.Institution_name, userHCISDetails.Institution_code).Scan(&instiDetails).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", insErr, insErr.Error())
			if !returnMessage.Data.IsSuccess {
				return returnMessage, userHCISDetails
			}
		}
	}

	userHCISDetails.Institution_id = instiDetails.Institution_id
	userHCISDetails.Institution_code = instiDetails.Institution_code
	userHCISDetails.Institution_name = instiDetails.Institution_name

	// marshal the response struct
	userHCISDetailsByte, marshalErr := json.Marshal(userHCISDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "", marshalErr, marshalErr.Error())
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// log the activity
	middleware.ActivityLogger(username, userHCISDetails.Institution_code, appCode, moduleName, funcName, "200", methodUsed, endpoint, reqBody, userHCISDetailsByte, "Success", "", nil)

	successResponse := response.ReturnModel{
		Data: response.DataModel{IsSuccess: true},
	}
	return successResponse, userHCISDetails
}
