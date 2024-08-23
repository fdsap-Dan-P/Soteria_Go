package registernewuser

import (
	"encoding/base64"
	"encoding/json"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"

	"github.com/go-resty/resty/v2"
)

func HcisInquiry(username, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) (response.ReturnModel, response.UserDetails) {
	userHCISDetails := response.UserDetails{}
	userHCISInfo := response.StaffInfoResponse{}
	instiDetails := response.InstitutionDetails{}

	funcName := "HCIS Inquiry"

	// Basic auth credentials
	hcis_username := "fdsap_apis"
	hcis_password := "P@ssword123"
	hcis_authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(hcis_username+":"+hcis_password))

	// Create a new Resty client
	client := resty.New()

	// Send the request
	resp, respErr := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", hcis_authHeader).
		SetBody(reqBody).
		Post("https://ua-uat.cardmri.com:8555/HCISLink/WEBAPI/ExternalService/ViewStaffInfo")

	if respErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "317", methodUsed, endpoint, reqBody, []byte(""), "Reading HCIS Response Failed", respErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	// Unmarshal the response body into the struct
	if unmarshallErr := json.Unmarshal(resp.Body(), &userHCISInfo); unmarshallErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "310", methodUsed, endpoint, reqBody, []byte(""), "", unmarshallErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if len(userHCISInfo.StaffInfo) == 0 {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "404", methodUsed, endpoint, reqBody, []byte(""), "Staff ID Not Found in HCIS", nil)
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
	if fetchErr := database.DBConn.Raw("SELECT * FROM public.institutions WHERE institution_name = ?", userHCISDetails.Institution_name).First(&instiDetails).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, userHCISDetails
		}
	}

	if instiDetails.Institution_id == 0 {
		userHCISDetails.Institution_code = GenerateInstitutionCode(userHCISDetails.Institution_name, username, instiCode, appCode, moduleName, methodUsed, endpoint)
		if insErr := database.DBConn.Raw("INSERT INTO public.institutions (institution_name, institution_code) VALUES (?, ?) RETURNING *", userHCISDetails.Institution_name, userHCISDetails.Institution_code).Scan(&instiDetails).Error; insErr != nil {
			returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", insErr)
			if !returnMessage.Data.IsSuccess {
				return returnMessage, userHCISDetails
			}
		}
	}

	userHCISDetails.Institution_id = instiDetails.Institution_id

	// marshal the response struct
	userHCISDetailsByte, marshalErr := json.Marshal(userHCISDetails)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "311", methodUsed, endpoint, reqBody, []byte(""), "", marshalErr)
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
