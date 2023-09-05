package util

import (
	"crypto/md5"
	"encoding/hex"
)

func HashPassword(password string) string {
	hasher := md5.New()

	// Write the string to the hasher
	hasher.Write([]byte(password))

	// Get the resulting hash value
	hashedBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
