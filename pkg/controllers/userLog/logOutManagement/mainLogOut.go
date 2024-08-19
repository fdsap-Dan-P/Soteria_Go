package logoutmanagement

import (
	"fmt"
	"soteria_go/pkg/controllers/middleware"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MainLogOut(c *fiber.Ctx) error {
	username := c.Params("username")
	sessionId := c.Params("session_id")
	userResponse := response.UserAccountResponse{}
	currentDateTime := middleware.GetDateTime()
	methodUsed := c.Method()
	endpoint := c.Path()
	funcName := "User Log"

	userActivity := "Logged Out"
	if !currentDateTime.Data.IsSuccess {
		return c.JSON(currentDateTime)
	}

	// get user id
	if fetchErr := database.DBConn.Raw("SELECT * FROM user_accounts WHERE username = ? OR staff_id = ?", username, username).Scan(&userResponse).Error; fetchErr != nil {
		returnMessage := middleware.ResponseData(username, funcName, " 302", methodUsed, endpoint, []byte(""), []byte(""), "", fetchErr)
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	if userResponse.User_id == 0 {
		returnMessage := middleware.ResponseData(username, funcName, " 404", methodUsed, endpoint, []byte(""), []byte(""), "User Not Found", fmt.Errorf(username))
		if !returnMessage.Data.IsSuccess {
			return c.JSON(returnMessage)
		}
	}

	// if no session
	if strings.TrimSpace(sessionId) == "" || strings.TrimSpace(sessionId) == "null" {
		isLogOut := WithoutSessionLogOut(username, funcName, methodUsed, endpoint)
		if !isLogOut.Data.IsSuccess {
			return c.JSON(isLogOut)
		}

	} else { // if there is session
		isLogOut := WithSessionLogOut(username, funcName, methodUsed, endpoint, sessionId)
		if !isLogOut.Data.IsSuccess {
			return c.JSON(isLogOut)
		}
	}

	// Audit Trails
	auditTrailLogs := middleware.AuditTrailGeneration(username, username, userActivity, "Logged In", "Logged Out", funcName, methodUsed, endpoint)
	if !auditTrailLogs.Data.IsSuccess {
		return c.JSON(auditTrailLogs)
	}

	middleware.ActivityLogger(username, "User Log", "202", methodUsed, endpoint, []byte(""), []byte(""), "Successfully Log Out", "", nil)
	return c.JSON(response.ResponseModel{
		RetCode: "202",
		Message: "Successfully Log Out",
	})
}
