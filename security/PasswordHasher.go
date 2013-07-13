package security

import (
	"crypto/sha256"
	"fmt"
)

func HashPassword(username string, password string) string {
	data := []byte(fmt.Sprintf("%v:%v", username, password))

	hash := sha256.New()
	result := hash.Sum(data)

	return fmt.Sprintf(`"%x"`, result)
}
