package middleware

import (
	"fmt"
	"os"

	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Username   string `json:"username"`
	Insti_code string `json:"insti_code"`
	jwt.StandardClaims
}

func GenerateToken(username, instiCode, appCode, moduleName, methodUsed, endpoint string) (string, error) {
	funcName := "Generate Token"

	configResponse := response.ConfigDetails{}
	if err := database.DBConn.Debug().Raw("SELECT * FROM parameters.system_config_params WHERE config_code = 'jwt' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&configResponse).Error; err != nil {
		return "", err
	}

	if configResponse.Config_value == "" {
		return "", fmt.Errorf("config value not found")
	}

	duration, err := strconv.Atoi(configResponse.Config_value)
	if err != nil {
		return "", err
	}

	expiration := time.Now()
	switch configResponse.Config_type {
	case "Second":
		expiration = expiration.Add(time.Second * time.Duration(duration))
	case "Minute":
		expiration = expiration.Add(time.Minute * time.Duration(duration))
	case "Hour":
		expiration = expiration.Add(time.Hour * time.Duration(duration))
	case "Day":
		expiration = expiration.AddDate(0, 0, duration)
	}

	claims := &Claims{
		Username:   username,
		Insti_code: instiCode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	// log the token generation
	ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(username), []byte(""), "Successful", "", nil)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// validate token
func ParseToken(tokenString, appCode, moduleName, methodUsed, endpoint string) response.ReturnModel {
	funcName := "Parse Token"

	// check if token string has value
	if strings.TrimSpace(tokenString) == "" {
		returnMessage := ResponseData("", "", appCode, moduleName, funcName, "401", methodUsed, endpoint, []byte(tokenString), []byte(""), "Token Missing", fmt.Errorf("token null"), nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		// Check if the error is due to token expiration
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				returnMessage := ResponseData("", "", appCode, moduleName, funcName, "110", methodUsed, endpoint, []byte(tokenString), []byte(""), "", ve, nil)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		}
		returnMessage := ResponseData("", "", appCode, moduleName, funcName, "301", methodUsed, endpoint, []byte(tokenString), []byte(""), "", err, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if !token.Valid {
		returnMessage := ResponseData("", "", appCode, moduleName, funcName, "104", methodUsed, endpoint, []byte(tokenString), []byte(""), "", fmt.Errorf("%s", tokenString), token)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// Token is valid, extract claims
	if claims, ok := token.Claims.(*Claims); ok {
		ActivityLogger("", "", appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(tokenString), []byte(string(claims.Username)), "Successful", "", nil)
		return response.ReturnModel{
			RetCode: "200",
			Message: string(claims.Username),
			Data: response.DataModel{
				Message:   string(claims.Insti_code),
				IsSuccess: true,
				Error:     nil,
			},
		}
	}

	returnMessage := ResponseData("", "", appCode, moduleName, funcName, "104", methodUsed, endpoint, []byte(tokenString), []byte(""), "", fmt.Errorf("%s", tokenString), token)
	if !returnMessage.Data.IsSuccess {
		return (returnMessage)
	}
	return returnMessage
}

func StoringUserToken(tokenString, username, staffId, instiCode, appCode, moduleName, methodUsed, endpoint string, reqBody []byte) response.ReturnModel {
	funcName := "Storing User Token"

	userTokenDetails := response.UserTokenDetails{}
	remark := response.DBFuncResponse{}

	// check if token string has value
	if strings.TrimSpace(tokenString) == "" {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "401", methodUsed, endpoint, []byte(tokenString), []byte(""), "New Generated Token Missing", fmt.Errorf("token null"), nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM logs.user_tokens WHERE (username = ? OR staff_id = ?) AND insti_code = ? AND app_code = ?", username, staffId, instiCode, appCode).Scan(&userTokenDetails).Error; fetchErr != nil {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "302", methodUsed, endpoint, reqBody, []byte(""), "", fetchErr, fetchErr.Error())
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if userTokenDetails.Token_id == 0 {
		if insErr := database.DBConn.Raw("SELECT logs.create_user_token(?, ?, ?, ?, ?) AS remark", username, staffId, tokenString, instiCode, appCode).Scan(&remark).Error; insErr != nil {
			returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", insErr, insErr.Error())
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		if remark.Remark != "Success" {
			returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "303", methodUsed, endpoint, reqBody, []byte(""), "", fmt.Errorf("%s", remark.Remark), remark)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	} else {
		if updErr := database.DBConn.Raw("SELECT logs.update_user_token(?, ?, ?, ?, ?) AS remark", tokenString, username, staffId, staffId, appCode).Scan(&remark).Error; updErr != nil {
			returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "304", methodUsed, endpoint, reqBody, []byte(""), "", updErr, updErr.Error())
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}

		if remark.Remark != "Success" {
			returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "304", methodUsed, endpoint, reqBody, []byte(""), "", fmt.Errorf("%s", remark.Remark), remark)
			if !returnMessage.Data.IsSuccess {
				return (returnMessage)
			}
		}
	}
	return response.ReturnModel{
		RetCode: "200",
		Message: "",
		Data: response.DataModel{
			Message:   "",
			IsSuccess: true,
			Error:     nil,
		},
	}
}
