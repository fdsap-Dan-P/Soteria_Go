package registernewuser

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
)

func HcisInquiry(staffId, username, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) (response.ReturnModel, response.UserDetails) {
	userHCISDetails := response.UserDetails{}
	userHCISInfo := response.StaffInfoResponse{}
	instiDetails := response.InstitutionDetails{}

	funcName := "HCIS Inquiry"

	if strings.TrimSpace(staffId) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, reqBody, []byte(""), "Staff Id Missing", nil, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Basic auth credentials [UAT]
	// url := "https://ua-uat.cardmri.com:8555/HCISLink/WEBAPI/ExternalService/ViewStaffInfo"
	// hcis_username := "unifiedauth"
	// hcis_password := "P@ssword123"

	// Basic auth credentials [PROD]
	url := "https://ua.cardmri.com:9125/HCISLink/WEBAPI/ExternalService/ViewStaffInfo"
	hcis_username := "cagabayfdsap"
	hcis_password := "ap!2024fdsap"
	hcis_authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(hcis_username+":"+hcis_password))

	hcis_reqBody := map[string]string{"StaffID": staffId}
	hcis_reqBodyByte, marshallErr := json.Marshal(hcis_reqBody)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "Marshalling HCIS Request Failed", marshallErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// // Create a new Resty client
	// client := resty.New()

	// // Send the request
	// resp, respErr := client.R().
	// 	SetHeader("Content-Type", "application/json").
	// 	SetHeader("Authorization", hcis_authHeader).
	// 	SetBody(hcis_reqBodyByte).
	// 	Post("https://ua-uat.cardmri.com:8555/HCISLink/WEBAPI/ExternalService/ViewStaffInfo")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	// Create the HTTP request and set headers
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(hcis_reqBodyByte))
	if reqErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "405", methodUsed, endpoint, hcis_reqBodyByte, []byte(""), "", reqErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage), userHCISDetails
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", hcis_authHeader)

	// Send the request
	resp, respErr := client.Do(req)
	if respErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "317", methodUsed, endpoint, hcis_reqBodyByte, []byte(""), "", respErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}
	defer resp.Body.Close()

	if respErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "317", methodUsed, endpoint, reqBody, []byte(""), "Reading HCIS Response Failed", respErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Check if the request was successful (status code 200)
	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "405", methodUsed, endpoint, reqBody, []byte(""), "Request Failed To HCIS", respErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "317", methodUsed, endpoint, (reqBody), []byte(""), "", err, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Unmarshal the response body into the struct
	if unmarshallErr := json.Unmarshal(body, &userHCISInfo); unmarshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "310", methodUsed, endpoint, reqBody, []byte(""), "", unmarshallErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if len(userHCISInfo.StaffInfo) == 0 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "Staff ID Not Found in HCIS", nil, userHCISDetails)
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
	if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_name = ?", userHCISDetails.Institution_name).Scan(&instiDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr, userHCISDetails)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if instiDetails.Institution_id == 0 {
		userHCISDetails.Institution_code = GenerateInstitutionCode(userHCISDetails.Institution_name, username, instiCode, appCode, moduleName, methodUsed, endpoint)
		if insErr := database.DBConn.Raw("INSERT INTO public.institutions (institution_name, institution_code) VALUES (?, ?) RETURNING *", userHCISDetails.Institution_name, userHCISDetails.Institution_code).Scan(&instiDetails).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", insErr, userHCISDetails)
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
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "", marshalErr, userHCISDetails)
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
