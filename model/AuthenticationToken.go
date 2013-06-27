package model

type AuthenticationToken string

func (token AuthenticationToken) String() string {
	return string(token)
}
