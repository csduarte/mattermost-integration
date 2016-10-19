package platform

import (
	"encoding/json"
	"io"
)

// User object
type User struct {
	ID                 string   `json:"id"`
	CreateAt           int64    `json:"create_at,omitempty"`
	UpdateAt           int64    `json:"update_at,omitempty"`
	DeleteAt           int64    `json:"delete_at"`
	Username           string   `json:"username"`
	Password           string   `json:"password,omitempty"`
	AuthData           *string  `json:"auth_data,omitempty"`
	AuthService        string   `json:"auth_service"`
	Email              string   `json:"email"`
	EmailVerified      bool     `json:"email_verified,omitempty"`
	Nickname           string   `json:"nickname"`
	FirstName          string   `json:"first_name"`
	LastName           string   `json:"last_name"`
	Roles              string   `json:"roles"`
	AllowMarketing     bool     `json:"allow_marketing,omitempty"`
	Props              []string `json:"props,omitempty"`
	NotifyProps        []string `json:"notify_props,omitempty"`
	LastPasswordUpdate int64    `json:"last_password_update,omitempty"`
	LastPictureUpdate  int64    `json:"last_picture_update,omitempty"`
	FailedAttempts     int      `json:"failed_attempts,omitempty"`
	Locale             string   `json:"locale"`
	MfaActive          bool     `json:"mfa_active,omitempty"`
	MfaSecret          string   `json:"mfa_secret,omitempty"`
	LastActivityAt     int64    `db:"-" json:"last_activity_at,omitempty"`
}

// UserFromJSON gets user from json
func UserFromJSON(data io.Reader) *User {
	decoder := json.NewDecoder(data)
	var user User
	err := decoder.Decode(&user)
	if err == nil {
		return &user
	}
	return nil
}
