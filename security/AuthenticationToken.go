package security

import (
	"github.com/nu7hatch/gouuid"
	"github.com/pjvds/httpcallback.io/model"
)

func NewAuthToken() model.AuthenticationToken {
	Log.Debug("Generating new AuthenticationToken")
	uuid, err := uuid.NewV4()
	if err != nil {
		Log.Fatalf("Error while generating new uuid V4: %s", err.Error())
	}

	return model.AuthenticationToken(uuid.String())
}
