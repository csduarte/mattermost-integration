package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// OutgoingWebhook structure
type OutgoingWebhook struct {
	ID           string   `json:"id"`
	Token        string   `json:"token"`
	CreateAt     int64    `json:"create_at"`
	UpdateAt     int64    `json:"update_at"`
	DeleteAt     int64    `json:"delete_at"`
	CreatorID    string   `json:"creator_id"`
	ChannelID    string   `json:"channel_id"`
	TeamID       string   `json:"team_id"`
	TriggerWords []string `json:"trigger_words"`
	TriggerWhen  int      `json:"trigger_when"`
	CallbackURLs []string `json:"callback_urls"`
	DisplayName  string   `json:"display_name"`
	Description  string   `json:"description"`
	ContentType  string   `json:"content_type"`
}

// OutgoingWebhookPayload structure
type OutgoingWebhookPayload struct {
	Token       string `json:"token"`
	TeamID      string `json:"team_id"`
	TeamDomain  string `json:"team_domain"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Timestamp   int64  `json:"timestamp"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	PostID      string `json:"post_id"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
}

// ToJSON marshal payload
func (o *OutgoingWebhookPayload) ToJSON() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(b)
}

// ToFormValues form value from struct
func (o *OutgoingWebhookPayload) ToFormValues() string {
	v := url.Values{}
	v.Set("token", o.Token)
	v.Set("team_id", o.TeamID)
	v.Set("team_domain", o.TeamDomain)
	v.Set("channel_id", o.ChannelID)
	v.Set("channel_name", o.ChannelName)
	v.Set("timestamp", strconv.FormatInt(o.Timestamp/1000, 10))
	v.Set("user_id", o.UserID)
	v.Set("user_name", o.UserName)
	v.Set("post_id", o.PostID)
	v.Set("text", o.Text)
	v.Set("trigger_word", o.TriggerWord)

	return v.Encode()
}

// ToJSON marshal struct
func (o *OutgoingWebhook) ToJSON() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(b)
}

// OutgoingWebhookFromJSON decode from json
func OutgoingWebhookFromJSON(data io.Reader) *OutgoingWebhook {
	decoder := json.NewDecoder(data)
	var o OutgoingWebhook
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	}
	return nil
}

// OutgoingWebhookPayloadFromJSON decode from json
func OutgoingWebhookPayloadFromJSON(data io.Reader) *OutgoingWebhookPayload {
	decoder := json.NewDecoder(data)
	var o OutgoingWebhookPayload
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	}
	return nil
}

// OutgoingWebhookPayloadFromForm decodes from url encoding
func OutgoingWebhookPayloadFromForm(data io.Reader) *OutgoingWebhookPayload {
	var o OutgoingWebhookPayload
	d, err := ioutil.ReadAll(data)
	if err != nil {
		return nil
	}
	f, err := url.ParseQuery(string(d))
	if err != nil {
		return nil
	}
	o.Token = f.Get("token")
	o.TeamID = f.Get("team_id")
	o.TeamDomain = f.Get("team_domain")
	o.ChannelID = f.Get("channel_id")
	o.ChannelName = f.Get("channel_name")
	// o.Timestamp = f.Get("timestamp")
	o.UserID = f.Get("user_id")
	o.UserName = f.Get("user_name")
	o.PostID = f.Get("post_id")
	o.Text = f.Get("text")
	o.TriggerWord = f.Get("trigger_word")
	return &o
}

// IsValid check valid webhook
func (o *OutgoingWebhook) IsValid() *ClientError {

	if len(o.ID) != 26 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.id.app_error", nil, "")
	}

	if len(o.Token) != 26 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.token.app_error", nil, "")
	}

	if o.CreateAt == 0 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.create_at.app_error", nil, "id="+o.ID)
	}

	if o.UpdateAt == 0 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.update_at.app_error", nil, "id="+o.ID)
	}

	if len(o.CreatorID) != 26 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.user_id.app_error", nil, "")
	}

	if len(o.ChannelID) != 0 && len(o.ChannelID) != 26 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.channel_id.app_error", nil, "")
	}

	if len(o.TeamID) != 26 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.team_id.app_error", nil, "")
	}

	if len(fmt.Sprintf("%s", o.TriggerWords)) > 1024 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.words.app_error", nil, "")
	}

	if len(o.TriggerWords) != 0 {
		for _, triggerWord := range o.TriggerWords {
			if len(triggerWord) == 0 {
				return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.trigger_words.app_error", nil, "")
			}
		}
	}

	if len(o.CallbackURLs) == 0 || len(fmt.Sprintf("%s", o.CallbackURLs)) > 1024 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.callback.app_error", nil, "")
	}

	for _, callback := range o.CallbackURLs {
		if !IsValidHTTPUrl(callback) {
			return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.url.app_error", nil, "")
		}
	}

	if len(o.DisplayName) > 64 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.display_name.app_error", nil, "")
	}

	if len(o.Description) > 128 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.description.app_error", nil, "")
	}

	if len(o.ContentType) > 128 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.content_type.app_error", nil, "")
	}

	if o.TriggerWhen > 1 {
		return NewClientError("OutgoingWebhook.IsValid", "model.outgoing_hook.is_valid.content_type.app_error", nil, "")
	}

	return nil
}

// func (o *OutgoingWebhook) PreSave() {
// 	if o.Id == "" {
// 		o.Id = NewId()
// 	}
//
// 	if o.Token == "" {
// 		o.Token = NewId()
// 	}
//
// 	o.CreateAt = GetMillis()
// 	o.UpdateAt = o.CreateAt
// }

// func (o *OutgoingWebhook) PreUpdate() {
// 	o.UpdateAt = GetMillis()
// }

// HasTriggerWord check hook for trigger words
func (o *OutgoingWebhook) HasTriggerWord(word string) bool {
	if len(o.TriggerWords) == 0 || len(word) == 0 {
		return false
	}

	for _, trigger := range o.TriggerWords {
		if trigger == word {
			return true
		}
	}

	return false
}

// TriggerWordStartsWith checks starts with trigger
func (o *OutgoingWebhook) TriggerWordStartsWith(word string) bool {
	if len(o.TriggerWords) == 0 || len(word) == 0 {
		return false
	}

	for _, trigger := range o.TriggerWords {
		if strings.HasPrefix(word, trigger) {
			return true
		}
	}

	return false
}

// func OutgoingWebhookListToJson(l []*OutgoingWebhook) string {
// 	b, err := json.Marshal(l)
// 	if err != nil {
// 		return ""
// 	} else {
// 		return string(b)
// 	}
// }
//
// func OutgoingWebhookListFromJson(data io.Reader) []*OutgoingWebhook {
// 	decoder := json.NewDecoder(data)
// 	var o []*OutgoingWebhook
// 	err := decoder.Decode(&o)
// 	if err == nil {
// 		return o
// 	} else {
// 		return nil
// 	}
// }
