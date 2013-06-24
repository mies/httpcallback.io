package security

import (
	"crypto/sha256"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
)

func HashPassword(username string, password string, creationDate time.Time) string {
	data := []byte(fmt.Sprintf("%s:%s:%s", username, password, creationDate))

	hash := sha256.New()
	result := hash.Sum(data)

	return fmt.Sprintf(`"%x"`, result)
}

type AuthenticationToken string

func NewAuthToken() AuthenticationToken {
	token, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return AuthenticationToken(token.String())
}

func (token AuthenticationToken) String() string {
	return fmt.Sprintf(`"%x"`, string(token.String()))
}

// func (token AuthenticationToken) MarshalJSON() ([]byte, error) {
// 	return []byte(fmt.Sprintf(`"%x"`, string(token.String()))), nil
// }

// // UnmarshalJSON turns *security.AuthenticationToken into a json.Unmarshaller.
// func (token *AuthenticationToken) UnmarshalJSON(data []byte) error {
// 	if len(data) != 53 || data[0] != '"' || data[52] != '"' {
// 		return errors.New(fmt.Sprintf("Invalid AuthenticationToken in JSON: %s", string(data)))
// 	}
// 	var buf [12]byte
// 	_, err := hex.Decode(buf[:], data[:])
// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Invalid AuthenicationToken in JSON: %s (%s)", string(data), err))
// 	}
// 	*token = AuthenticationToken(string(data[:]))
// 	return nil
// }
