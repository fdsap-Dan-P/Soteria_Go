package validations

import (
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
)

func ValidateInstiCode(instiCode, username, funcName, methodUsed, endpoint string) response.ReturnModel {
	instiResponse := response.InstitutionDetails{}
	successResponse := response.ReturnModel{}

	if strings.TrimSpace(instiCode) == "" || strings.TrimSpace(instiCode) == "null" {
		returnMessage := middleware.ResponseData(username, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "Institution Code Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// check if application code exists
	if fetchErr := database.DBConn.Raw("SELECT * FROM offices_mapping.institutions WHERE institution_code = ?", instiCode).Scan(&instiResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if instiResponse.Institution_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, "404", methodUsed, endpoint, []byte(""), []byte(""), "Institution Code Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	successResponse.Data.IsSuccess = true
	return (successResponse)
}
