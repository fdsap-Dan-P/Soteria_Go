package setuserpassword

import (
	"fmt"
	"regexp"
	"soteria_go/pkg/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/hash"
	"strconv"
	"strings"
)

func PasswordValidation(password, instiCode, appCode, username, moduleName, methodUsed, endpoint string) response.ReturnModel {
	funcName := "Password Validation"
	// check the password:
	// length
	// at least one capital letter
	// at least one small letter
	// at least one number
	// at least one special character
	// - - - - - - - - P A S S W O R D    L E N G T H    V A L I D A T I O N - - - - - - - -//
	passMin := response.ConfigDetails{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config_params WHERE config_code = 'pass_min' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&passMin).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if passMin.Config_id == 0 || strings.TrimSpace(passMin.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "121", methodUsed, endpoint, []byte(""), []byte(""), "No Config Password Mininum Set", nil, passMin)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// parse the password minimum into int
	passMinValue, parsErr := strconv.Atoi(passMin.Config_value)
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Password Minimum Failed", parsErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if len(password) < passMinValue {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "103", methodUsed, endpoint, []byte(""), []byte(""), "Invalid Password", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// - - - - - - - - P A S S W O R D    C H A R A C T E R S    V A L I D A T I O N - - - - - - - -//
	passCharValidation := []response.ConfigDetails{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config_params WHERE config_value = 'true' AND config_code IN('pass_lowc','pass_upc','pass_num','pass_sym')").Scan(&passCharValidation).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// Regular expressions to check for lowercase, uppercase, symbol, and number
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	symbolRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	numberRegex := regexp.MustCompile(`[0-9]`)

	// Check each configuration against the password
	for _, configResponse := range passCharValidation {
		switch configResponse.Config_code {
		case "pass_lowc":
			fmt.Println("pass_lowc: true")
			if !lowercaseRegex.MatchString(password) {
				fmt.Println(password + "has no LOWERCASE")
				returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "103", methodUsed, endpoint, []byte(""), []byte(""), "Invalid Password | No Lowercase", nil, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		case "pass_upc":
			fmt.Println("pass_upc: true")
			if !uppercaseRegex.MatchString(password) {
				fmt.Println(password + "has no UPPERCASE")
				returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "103", methodUsed, endpoint, []byte(""), []byte(""), "Invalid Password | No Uppercase", nil, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		case "pass_sym":
			fmt.Println("pass_sym: true")
			if !symbolRegex.MatchString(password) {
				fmt.Println(password + "has no SPECIAL CHARACTERS")
				returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "103", methodUsed, endpoint, []byte(""), []byte(""), "Invalid Password | No Special Characters", nil, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		case "pass_num":
			fmt.Println("pass_num: true")
			if !numberRegex.MatchString(password) {
				fmt.Println(password + "has no NUMBERS")
				returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "103", methodUsed, endpoint, []byte(""), []byte(""), "Invalid Password | No Numbers", nil, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		}
	}

	successResp := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Validated New Password", nil, nil)
	if !successResp.Data.IsSuccess {
		return (successResp)
	}

	return successResp
}

func PasswordReuseValidation(password string, instiCode string, appCode string, username string, moduleName string, methodUsed string, endpoint string, userId int) response.ReturnModel {
	prevPasswords := []response.LastPasswordUsed{}
	minPassReuse := response.ConfigDetails{}

	funcName := "Password Reuse Validation"

	// get the minimum password reuse
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config_params WHERE config_code = 'pass_reuse' AND config_insti_code = ? AND config_app_code = ?").Scan(&minPassReuse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if minPassReuse.Config_id == 0 || strings.TrimSpace(minPassReuse.Config_value) == "" {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "121", methodUsed, endpoint, []byte(""), []byte(""), "No Minimum Password Reuse Config Set", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// parse the password minimum reuse value into int
	minPassReuseInt, parsErr := strconv.Atoi(strings.TrimSpace(minPassReuse.Config_value))
	if parsErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Minimum Password Reuse Value Failed", parsErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// get the last password used
	if fetchErr := database.DBConn.Raw("SELECT * FROM password_reuse(?, ?)", userId, minPassReuseInt).Scan(&prevPasswords).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	newHashedPassword := hash.SHA256(password)

	// check if new password matches to previous password
	for _, previouseUserPassword := range prevPasswords {
		if newHashedPassword == previouseUserPassword.Password_hash {
			returnMessage := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "113", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	}

	succesResp := middleware.ResponseData(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(""), []byte(""), "", nil, nil)
	if !succesResp.Data.IsSuccess {
		return (succesResp)
	}

	return (succesResp)
}
