package models

type User struct {
	ID        string
	Email     string
	Username  string
	Name      string
	Bio       string
	AvatarURL string
	LastSeen  string
	PassHash  []byte
}
