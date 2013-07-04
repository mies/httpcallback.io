package security

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func HashPassword(username string, password string, creationDate time.Time) string {
	data := []byte(fmt.Sprintf("%s:%s:%s", username, password, creationDate))

	hash := sha256.New()
	result := hash.Sum(data)

	return fmt.Sprintf(`"%x"`, result)
}
