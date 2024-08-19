package middleware

import (
	"fmt"
	"os"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Staff_id string `json:"staff_id"`
	jwt.StandardClaims
}

func GenerateToken(userIdentity, instiCode, appCode string) (string, error) {
	configResponse := response.SystemConfigurationResponse{}
	userAccount := response.UserAccountResponse{}
	if err := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'jwt' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&configResponse).Error; err != nil {
		return "", err
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

	// get user unique idengtities
	if fetchErr := database.DBConn.Raw("SELECT * FROM parameters.system_config WHERE config_code = 'jwt' AND config_insti_code = ? AND config_app_code = ?", instiCode, appCode).Scan(&configResponse).Error; fetchErr != nil {
		return "", fetchErr
	}

	if userAccount.User_id == 0 {
		return "User Not Found", fmt.Errorf(userIdentity)
	}

	claims := &Claims{
		User_id:  userAccount.User_id,
		Username: userAccount.Username,
		Staff_id: userAccount.Staff_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// validate token
func ParseToken(tokenString, funcName, methodUsed, endpoint string) response.ReturnModel {
	if strings.TrimSpace(tokenString) == "" {
		returnMessage := ResponseData("", funcName, "401", methodUsed, endpoint, []byte(tokenString), []byte(""), "Session Id Not Found", fmt.Errorf("token null"))
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
				returnMessage := ResponseData("", funcName, "110", methodUsed, endpoint, []byte(tokenString), []byte(""), "", ve)
				if !returnMessage.Data.IsSuccess {
					return (returnMessage)
				}
			}
		}
		returnMessage := ResponseData("", funcName, "301", methodUsed, endpoint, []byte(tokenString), []byte(""), "", err)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if !token.Valid {
		returnMessage := ResponseData("", funcName, "104", methodUsed, endpoint, []byte(tokenString), []byte(""), "", fmt.Errorf(tokenString))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// Token is valid, extract claims
	if claims, ok := token.Claims.(*Claims); ok {
		ActivityLogger("", funcName, "200", methodUsed, endpoint, []byte(tokenString), []byte(string(claims.Username)), "Successful", "", nil)
		return response.ReturnModel{
			RetCode: "200",
			Message: "Successful",
			Data: errors.ErrorModel{
				Message:   string(claims.Username),
				User_id:   claims.User_id,
				Username:  claims.Username,
				Staff_id:  claims.Staff_id,
				IsSuccess: true,
				Error:     nil,
			},
		}
	}

	returnMessage := ResponseData("", funcName, "104", methodUsed, endpoint, []byte(tokenString), []byte(""), "", fmt.Errorf(tokenString))
	if !returnMessage.Data.IsSuccess {
		return (returnMessage)
	}
	return returnMessage
}
