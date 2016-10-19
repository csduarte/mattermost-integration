package platform

import (
	"encoding/json"
	"io"
)

// Channel struct
type Channel struct {
	ID            string `json:"id"`
	CreateAt      int64  `json:"create_at"`
	UpdateAt      int64  `json:"update_at"`
	DeleteAt      int64  `json:"delete_at"`
	TeamID        string `json:"team_id"`
	Type          string `json:"type"`
	DisplayName   string `json:"display_name"`
	Name          string `json:"name"`
	Header        string `json:"header"`
	Purpose       string `json:"purpose"`
	LastPostAt    int64  `json:"last_post_at"`
	TotalMsgCount int64  `json:"total_msg_count"`
	ExtraUpdateAt int64  `json:"extra_update_at"`
	CreatorID     string `json:"creator_id"`
}

// ChannelList struct
type ChannelList struct {
	Channels []*Channel `json:"channels"`
}

// ChannelListFromJSON decodes channellist response
func ChannelListFromJSON(data io.Reader) *ChannelList {
	decoder := json.NewDecoder(data)
	var o ChannelList
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	}
	return nil
}
