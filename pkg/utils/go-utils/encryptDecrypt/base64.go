package encryptDecrypt

import (
	"encoding/base64"
	"log"
)

func EncodeBase64(plaintext string) string {
	// Encrypt the text using Base64 encoding
	encrypted := base64.StdEncoding.EncodeToString([]byte(plaintext))

	return encrypted
}

func DecodeBase64(encrypted string) string {
	// Decrypt the text using Base64 decoding
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		log.Println("Error decoding Base64:", err)
		return err.Error()
	}

	decodedString := string(decoded)

	return decodedString
}
