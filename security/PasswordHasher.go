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

type AuthenticationToken struct {
	token uuid.UUID
}

func NewAuthToken() AuthenticationToken {
	Log.Debug("Generating new AuthenticationToken")
	uuid, err := uuid.NewV4()
	if err != nil {
		Log.Fatalf("Error while generating new uuid V4: %s", err.Error())
	}

	return AuthenticationToken{
		token: *uuid,
	}
}

func (token AuthenticationToken) String() string {
	return token.token.String()
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
