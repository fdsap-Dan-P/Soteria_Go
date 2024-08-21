package registernewuser

import (
	"soteria_go/pkg/middleware"
	"strings"
)

func GenerateInstitutionCode(institutionName, username, instiCode, appCode, moduleName, methodUsed, endpoint string) string {
	funcName := "Generate Institution Code"
	// Convert the name to uppercase
	institutionName = strings.ToUpper(institutionName)

	// Split the name into words
	words := strings.Fields(institutionName)

	// Generate the code
	var generatedInstiCode string
	for _, word := range words {
		generatedInstiCode += string(word[0]) // Take the first letter of each word
	}

	middleware.ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(instiCode), []byte(generatedInstiCode), "Success", "", nil)

	return generatedInstiCode
}
