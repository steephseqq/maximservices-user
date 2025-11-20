package models

type User struct {
	UUID      string
	Email     string
	Username  string
	Name      string
	Bio       string
	AvatarURL string
	LastSeen  string
	PassHash  []byte
}
