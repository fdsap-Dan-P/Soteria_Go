package middleware

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	otpLength       = 6
	otpExpiryPeriod = 5 * time.Minute
)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	return strconv.Itoa(otp)
}

func OTPGeneration() (string, string) {
	otp := generateOTP()
	creationTime := time.Now()

	// Validate the OTP and check expiration
	isValid, isExpired := validateOTP(otp, creationTime)
	fmt.Println("Is OTP valid?", isValid)
	fmt.Println("Is OTP expired?", isExpired)
	createdOTP := creationTime.Format("2006-01-02 15:04:05")
	return otp, createdOTP
}

func validateOTP(otp string, creationTime time.Time) (bool, bool) {
	if len(otp) != otpLength {
		return false, false
	}

	otpInt, err := strconv.Atoi(otp)
	if err != nil {
		return false, false
	}

	if time.Since(creationTime) > otpExpiryPeriod {
		return false, true // OTP is invalid due to expiration
	}

	return otpInt >= 100000 && otpInt <= 999999, false
}
