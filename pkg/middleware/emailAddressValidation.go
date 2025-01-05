package middleware

import "regexp"

func ValidateEmail(email string) bool {
	// Regular expression for validating email addresses
	// This regex is a simplified version and may not cover all edge cases
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$` // only check if .com

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the email matches the pattern
	return regex.MatchString(email)
}
