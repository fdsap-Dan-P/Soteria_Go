package middleware

import (
	"fmt"
	"soteria_go/pkg/models/errors"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SessionIdChecker(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	funcName := "Report Module"
	methodUsed := c.Method()
	endpoint := c.Path()

	if strings.TrimSpace(uuid) == "" || strings.TrimSpace(uuid) == "null" {
		returnMessage := ResponseData("", funcName, "115", methodUsed, endpoint, []byte(uuid), []byte(""), "", fmt.Errorf("null uuid"))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}
	sessionDetails := response.SessionLogs{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM sessions WHERE session_id = ?", uuid).Scan(&sessionDetails).Error; fetchErr != nil {
		returnMessage := ResponseData("", funcName, "302", methodUsed, endpoint, []byte(uuid), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if sessionDetails.User_id == 0 {
		returnMessage := ResponseData("", funcName, "404", methodUsed, endpoint, []byte(uuid), []byte(""), "Session Id Not Found", fmt.Errorf(uuid))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Successful",
		Data:    sessionDetails,
	})
}

func CreateAndStoreSessionId(username, funcName, methodUsed, endpoint, currentDateTime, userIP string, userID int, reqByte []byte) response.ReturnModel {
	remark := response.DBFuncResponse{}
	// generate session id
	uuid := uuid.New()

	if deletErr := database.DBConn.Raw("SELECT delete_from_sessions(?) AS remark", userID).Scan(&remark).Error; deletErr != nil {
		returnMessage := ResponseData(username, funcName, "314", methodUsed, endpoint, (reqByte), []byte(""), "", deletErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	if insErr := database.DBConn.Raw("SELECT add_to_session (?, ?, ?, ?, ?) AS remark", uuid.String(), userID, currentDateTime, userIP, true).Scan(&remark).Error; insErr != nil {
		returnMessage := ResponseData(username, funcName, "303", methodUsed, endpoint, (reqByte), []byte(""), "", insErr)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}
	if remark.Remark != "Success" {
		returnMessage := ResponseData(username, funcName, "303", methodUsed, endpoint, (reqByte), []byte(""), "", fmt.Errorf(remark.Remark))
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: errors.ErrorModel{
			Message:   uuid.String(),
			IsSuccess: true,
			Error:     nil,
		},
	}
}
