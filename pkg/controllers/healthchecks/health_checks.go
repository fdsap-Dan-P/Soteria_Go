package healthchecks

import (
	"soteria_go/pkg/models/response"

	"github.com/gofiber/fiber/v2"
)

func checkHealth() response.ResponseModel {
	return response.ResponseModel{
		RetCode: "100",
		Message: "Request success!",
		Data: response.DataModel{
			Message:   "Service is available!",
			IsSuccess: true,
			Error:     nil,
		},
	}
}

func checkHealthB() response.ResponseModel {
	return response.ResponseModel{
		RetCode: "100",
		Message: "Request success B!",
		Data: response.DataModel{
			Message:   "Service is available!",
			IsSuccess: true,
			Error:     nil,
		},
	}
}

func CheckServiceHealth(c *fiber.Ctx) error {
	health := checkHealth()
	healthResponse := response.DataModel{}
	healthResponse = health.Data.(response.DataModel)
	if !healthResponse.IsSuccess {
		return c.JSON(health)
	}
	return c.JSON(health)
}

func CheckServiceHealthB(c *fiber.Ctx) error {
	health := checkHealthB()
	healthResponse := response.DataModel{}
	healthResponse = health.Data.(response.DataModel)
	if !healthResponse.IsSuccess {
		return c.JSON(health)
	}
	return c.JSON(health)
}
