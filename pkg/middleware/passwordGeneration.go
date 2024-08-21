package middleware

import (
	"math/rand"
	"time"
)

func PasswordGeneration() string {
	// Define the characters that can be used in the password (excluding "P@ssw0rd")
	var characters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a byte slice to store the password characters
	password := make([]rune, 20)

	// Set the initial characters to "P@ssw0rd"
	initial := []rune("T3mpP@ssw0rd-")
	copy(password, initial)

	// Generate random characters for the last 4 characters of the password
	for i := 13; i < 20; i++ {
		// Generate a random index to select a character from the character set
		password[i] = characters[rand.Intn(len(characters))]
	}

	// no need logs no error, success was log in main function
	return string(password)
}
