package memberVerification

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
)

func GetMemberSavingAcountDetailsFromDataMart(cid, instiCode, fullName, appCode, moduleName, methodUsed, endpoint, apiKey string) (response.ReturnModel, response.MemberResponse) {
	dmUserSavings := response.MemberSavingResponse{}

	funcName := "Get Member Saving Account Details From Data Mart"

	// make request to data mart to  get user's saving account
	dmMemberSavingRequestBody := request.MemberVerificationRequest{
		Cid:              cid,
		Institution_code: instiCode,
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
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "311", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "Marshalling Request Body Failed", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}

	// Create the HTTP request and set headers
	dmMemberSavingReq, dmMemberSavingReqErr := http.NewRequest("POST", dmMemberSavingUrl, bytes.NewBuffer(dmMemberSavingRequestByte))
	if dmMemberSavingReqErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "318", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "", dmMemberSavingReqErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}
	dmMemberSavingReq.Header.Set("Content-Type", "application/json")
	dmMemberSavingReq.Header.Set("X-API-Key", apiKey)

	fmt.Println(dmMemberSavingReq)

	// Send the request
	dmMemberSavingResp, dmMemberSavingRespErr := dmMemberSavingClient.Do(dmMemberSavingReq)
	if dmMemberSavingRespErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "317", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "", dmMemberSavingRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}
	defer dmMemberSavingResp.Body.Close()

	if dmMemberSavingResp.Status != "200 OK" {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "405", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "", dmMemberSavingRespErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}

	dmMemberSavingBody, dmMemberSavingErr := ioutil.ReadAll(dmMemberSavingResp.Body)
	if dmMemberSavingErr != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "317", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "Reading Data Mart Response Failed", dmMemberSavingErr, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}

	// soteriaResponse_v2 := map[string]string{}
	if err := json.Unmarshal(dmMemberSavingBody, &dmUserSavings); err != nil {
		returnMessage := middleware.ResponseData(fullName, "", appCode, moduleName, funcName, "310", methodUsed, endpoint, dmMemberSavingRequestByte, []byte(""), "", err, nil)
		if !returnMessage.Data.IsSuccess {
			return returnMessage, dmUserSavings
		}
	}

	successMessage := response.ReturnModel{Data: response.DataModel{IsSuccess: true}}

	return successMessage, dmUserSavings
}
