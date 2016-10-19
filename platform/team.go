package platform

import (
	"encoding/json"
	"io"
)

// Team structure
type Team struct {
	ID              string `json:"id"`
	CreateAt        int64  `json:"create_at"`
	UpdateAt        int64  `json:"update_at"`
	DeleteAt        int64  `json:"delete_at"`
	DisplayName     string `json:"display_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Type            string `json:"type"`
	CompanyName     string `json:"company_name"`
	AllowedDomains  string `json:"allowed_domains"`
	InviteID        string `json:"invite_id"`
	AllowOpenInvite bool   `json:"allow_open_invite"`
}

// TeamMapFromJSON fetches map of all teams
func TeamMapFromJSON(data io.Reader) map[string]*Team {
	decoder := json.NewDecoder(data)
	var teams map[string]*Team
	err := decoder.Decode(&teams)
	if err == nil {
		return teams
	}
	return nil
}
