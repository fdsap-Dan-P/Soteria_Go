package validations

import (
	"regexp"
)

func ValidateStaffID(staffID string) bool {
	// Define the regular expression pattern for staff ID
	pattern := `^\d{6}-\d{5}$`

	// Compile the regular expression
	regexp := regexp.MustCompile(pattern)

	// Check if the staff ID matches the pattern
	return regexp.MatchString(staffID)
}
