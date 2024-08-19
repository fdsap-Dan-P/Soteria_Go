package changeuserpassword

import (
	"encoding/json"
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/request"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsernameForgotPassword(c *fiber.Ctx) error {
	user := request.ForgotPassword{}
	userInfoResponse := response.UserInfo{}
	remark := response.DBFuncResponse{}

	funcName := "Forgot Password"
	methodUsed := c.Method()
	endpoint := c.Path()
	userIP := c.IP()

	currentDateTime := middleware.GetDateTime().Data.Message
	if parsErr := c.BodyParser(&user); parsErr != nil {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "301", methodUsed, endpoint, []byte(""), []byte(""), "Parsing Username Failed", parsErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if strings.TrimSpace(user.Username) == "" {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "401", methodUsed, endpoint, []byte(""), []byte(""), "Username or Employee ID Input Missing", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// marshal request body
	requestBodyBytes, marshallErr := json.Marshal(user)
	if marshallErr != nil {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "311", methodUsed, endpoint, []byte(""), []byte(""), "", marshallErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if fetchErr := database.DBConn.Raw("SELECT * FROM user_info WHERE username = ? OR staff_id = ?", user.Username, user.Username).Scan(&userInfoResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "302", methodUsed, endpoint, requestBodyBytes, []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userInfoResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "404", methodUsed, endpoint, requestBodyBytes, []byte(""), "Username or Employee ID Not Found", nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userInfoResponse.Status != "Active" {
		errMessage := fmt.Sprintf("%s User Account", userInfoResponse.Status)
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "123", methodUsed, endpoint, (requestBodyBytes), []byte(""), errMessage, nil)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	uuid := uuid.New()
	if insErr := database.DBConn.Raw("SELECT add_to_session (?, ?, ?, ?, ?) AS remark", uuid.String(), userInfoResponse.User_id, currentDateTime, userIP, true).Scan(&remark).Error; insErr != nil {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "303", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", insErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := middleware.ResponseData(userInfoResponse.Username, funcName, "303", methodUsed, endpoint, (requestBodyBytes), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	// insert session id
	middleware.ActivityLogger(userInfoResponse.Username, funcName, "200", methodUsed, endpoint, requestBodyBytes, []byte(""), "Successful", "", nil)
	return c.JSON(response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:    userInfoResponse.Username,
			User_id:    userInfoResponse.User_id,
			Session_id: uuid.String(),
			Is_active:  true, // session id status
		},
	})
}
