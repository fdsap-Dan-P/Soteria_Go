package securitymanagement

import (
	"os"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"soteria_go/pkg/utils/go-utils/encryptDecrypt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func EncryptApiKey(c *fiber.Ctx) error {
	appDetailsList := []response.ApplicationDetails{}
	insertedApp := response.ApplicationDetails{}

	if fetchErr := database.DBConn.Raw("SELECT * FROM public.applications").Scan(&appDetailsList).Error; fetchErr != nil {
		return c.JSON(response.ReturnModel{
			RetCode: "302",
			Message: "Internal Server Error",
			Data: response.DataModel{
				Message:   "Fetching Data Failed",
				IsSuccess: false,
				Error:     fetchErr,
			},
		})
	}

	if len(appDetailsList) == 0 {
		return c.JSON(response.ReturnModel{
			RetCode: "106",
			Message: "Validation Failed",
			Data: response.DataModel{
				Message:   "No Data Available",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	// get the secret key
	secretKey := os.Getenv("SECRET_KEY")
	if strings.TrimSpace(secretKey) == "" {
		return c.JSON(response.ReturnModel{
			RetCode: "404",
			Message: "Bad Request",
			Data: response.DataModel{
				Message:   "Secret Key Not Found in Environment",
				IsSuccess: false,
				Error:     nil,
			},
		})
	}

	for _, appDetails := range appDetailsList {
		// encrypt the user api key
		encryptedApiKey, encryptErr := encryptDecrypt.EncryptWithSecretKey(appDetails.Api_key, secretKey)
		if encryptErr != nil {
			return c.JSON(response.ReturnModel{
				RetCode: "318",
				Message: "Internal Server Error",
				Data: response.DataModel{
					Message:   "Encrypting Api Key Failed",
					IsSuccess: false,
					Error:     encryptErr,
				},
			})
		}

		if updatErr := database.DBConn.Raw("UPDATE public.applications SET api_key = ? WHERE application_code = ?", encryptedApiKey, appDetails.Application_code).Scan(&insertedApp).Error; updatErr != nil {
			return c.JSON(response.ReturnModel{
				RetCode: "304",
				Message: "Internal Server Error",
				Data: response.DataModel{
					Message:   "Updating Data Failed",
					IsSuccess: false,
					Error:     updatErr,
				},
			})
		}
	}

	return c.JSON(response.ResponseModel{
		RetCode: "204",
		Message: "Successfully Updated",
	})
}
