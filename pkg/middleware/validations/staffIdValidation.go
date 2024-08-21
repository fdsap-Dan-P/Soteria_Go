package validations

import (
	"fmt"
	"regexp"
	"soteria_go/pkg/middleware"
)

func StaffIdValidation(staffID, moduleName, methodUsed, endpoint string) bool {
	funcName := "Validate Staff ID Format"
	// Define the regular expression pattern for staff ID
	pattern := `^\d{6}-\d{5}$`

	// Compile the regular expression
	regexp := regexp.MustCompile(pattern)

	// log the activity
	respBody := fmt.Sprintf("Is Staff ID Valid: %v", regexp.MatchString(staffID))
	middleware.ActivityLogger(staffID, "", "", moduleName, funcName, "200", methodUsed, endpoint, []byte(staffID), []byte(respBody), "Successful", respBody, nil)
	// Check if the staff ID matches the pattern
	return regexp.MatchString(staffID)
}
