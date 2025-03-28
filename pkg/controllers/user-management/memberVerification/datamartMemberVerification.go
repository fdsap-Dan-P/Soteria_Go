package memberVerification

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
)

func VerifyMemberFromDataMart(fullName, appCode, moduleName, methodUsed, endpoint, apiKey string, userRequest request.MemberVerificationRequest) (response.ReturnModel, response.MemberResponse) {
	dmUserDetails := response.MemberResponse{}

	funcName := "Verify Member From Data Mart"

	dmMemberVerificationUrl := os.Getenv("DATA_MART_HOST") + "/users/verify-member"

	dmMemberVerificationtr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	dmMemberVerificationClient := &http.Client{Transport: dmMemberVerificationtr}

	// marshal the request body for data mart
	dmMemberVerificationRequestByte, marshalErr := json.Marshal(userRequest)
	if marshalErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "Marshalling Request Body Failed", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}

	// Create the HTTP request and set headers
	dmMemberVerificationReq, dmMemberVerificationReqErr := http.NewRequest("POST", dmMemberVerificationUrl, bytes.NewBuffer(dmMemberVerificationRequestByte))
	if dmMemberVerificationReqErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "318", methodUsed, endpoint, dmMemberVerificationRequestByte, []byte(""), "", dmMemberVerificationReqErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}
	dmMemberVerificationReq.Header.Set("Content-Type", "application/json")
	dmMemberVerificationReq.Header.Set("X-API-Key", apiKey)

	// Send the request
	dmMemberVerificationResp, dmMemberVerificationRespErr := dmMemberVerificationClient.Do(dmMemberVerificationReq)
	if dmMemberVerificationRespErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "317", methodUsed, endpoint, dmMemberVerificationRequestByte, []byte(""), "", dmMemberVerificationRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}
	defer dmMemberVerificationResp.Body.Close()

	if dmMemberVerificationResp.Status != "200 OK" {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "405", methodUsed, endpoint, dmMemberVerificationRequestByte, []byte(""), "", dmMemberVerificationRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}

	dmMemberVerificationBody, dmMemberVerificationErr := ioutil.ReadAll(dmMemberVerificationResp.Body)
	if dmMemberVerificationErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "317", methodUsed, endpoint, dmMemberVerificationRequestByte, []byte(""), "Reading Data Mart Response Failed", dmMemberVerificationErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}

	// soteriaResponse_v2 := map[string]string{}
	if err := json.Unmarshal(dmMemberVerificationBody, &dmUserDetails); err != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "310", methodUsed, endpoint, dmMemberVerificationRequestByte, []byte(""), "", err, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserDetails
		}
	}

	successMessage := response.ReturnModel{Data: response.DataModel{IsSuccess: true}}

	return successMessage, dmUserDetails
}
