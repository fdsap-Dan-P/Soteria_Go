package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
