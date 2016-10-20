package server

import (
	"fmt"

	"github.com/csduarte/mattermost-integration/platform"
)

// MattermostMessage A convience holder to represent mattermost Incoming Webbhook
type MattermostMessage struct {
	platform.IncomingWebhookRequest
	integration *Integration
}

// NewMattermostMessage initializes a new message
func NewMattermostMessage(channel string) *MattermostMessage {
	m := MattermostMessage{}
	m.ChannelName = channel
	return &m
}

// AddIconURL will simple add any string as the IconURL
func (m *MattermostMessage) AddIconURL(url string) {
	m.IconURL = url
}

// AddImageURL will add a simple string
func (m *MattermostMessage) AddImageURL(url string) {
	// m.ensureAttachments()
	// m.Attachments.(map[string][]interface{})["image_url"] = url
	m.Attachments = []map[string]interface{}{
		0: {
			"image_url": url,
		},
	}
}

// SetUsername sets username override, if allowed on server
func (m *MattermostMessage) SetUsername(name string) {
	m.Username = name
}

// SetMessage sets the text value
func (m *MattermostMessage) SetMessage(msg string) {
	m.Text = msg
}

// Send message to associated integration
func (m *MattermostMessage) Send() error {
	if m.integration == nil {
		return fmt.Errorf("Could not send message - missing integration for username %q\n", m.Username)
	}
	if m.integration.ToMM != nil && m.integration.client != nil {
		in := m.integration
		data := m.IncomingWebhookRequest.ToJSON()
		fmt.Println(data)
		_, err := in.client.PostToWebhook(in.ToMM.ID, data)
		return fmt.Errorf(err.Message)
	}
	return fmt.Errorf("Could not send message - missing integration for %v\n", m.integration.Name)
}

// ensureAttachments make sure that the attachments has been set
func (m *MattermostMessage) ensureAttachments() {
	if m.Attachments == nil {
		m.Attachments = make(map[string][]interface{})
	}
}
