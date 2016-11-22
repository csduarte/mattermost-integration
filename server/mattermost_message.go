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

// SetIconURL will simple add any string as the IconURL
func (m *MattermostMessage) SetIconURL(url string) {
	m.IconURL = url
}

// SetUsername sets username override, if allowed on server
func (m *MattermostMessage) SetUsername(name string) {
	m.Username = name
}

// SetMessage sets the text value
func (m *MattermostMessage) SetMessage(msg string) {
	m.Text = msg
}

// AttachmentAuthorIcon sets author icon
func (m *MattermostMessage) AttachmentAuthorIcon(url string) {
	a := m.ensureAttachments()
	a.AuthorIcon = url
}

// AttachmentAuthorLink sets author link
func (m *MattermostMessage) AttachmentAuthorLink(url string) {
	a := m.ensureAttachments()
	a.AuthorLink = url
}

// AttachmentAuthorName sets author name
func (m *MattermostMessage) AttachmentAuthorName(text string) {
	a := m.ensureAttachments()
	a.AuthorName = text
}

// AttachmentColor sets attachment color (ex: #333333)
func (m *MattermostMessage) AttachmentColor(cssHexColor string) {
	a := m.ensureAttachments()
	a.Color = cssHexColor
}

// AttachmentFallback sets attachment fallback text
func (m *MattermostMessage) AttachmentFallback(text string) {
	a := m.ensureAttachments()
	a.Fallback = text
}

// AttachmentPretext sets attachment Pretext
func (m *MattermostMessage) AttachmentPretext(text string) {
	a := m.ensureAttachments()
	a.Pretext = text
}

// AttachmentText sets attachment text
func (m *MattermostMessage) AttachmentText(text string) {
	a := m.ensureAttachments()
	a.Text = text
}

// AttachmentTitle sets attachment title
func (m *MattermostMessage) AttachmentTitle(text string) {
	a := m.ensureAttachments()
	a.Title = text
}

// AttachmentTitleLink sets attachment title link
func (m *MattermostMessage) AttachmentTitleLink(url string) {
	a := m.ensureAttachments()
	a.TitleLink = url
}

// AttachmentImageURL will add a image url attachment
func (m *MattermostMessage) AttachmentImageURL(url string) {
	a := m.ensureAttachments()
	a.ImageURL = url
}

// AttachmentAddField adds a field to the message Attachment
func (m *MattermostMessage) AttachmentAddField(title, text string, short bool) *platform.AttachmentField {
	a := m.ensureFields()
	f := platform.AttachmentField{}
	f.Short = short
	f.Title = title
	f.Value = text
	a.Fields = append(a.Fields, &f)
	return &f
}

// ensureAttachments make sure that the attachments has been set
func (m *MattermostMessage) ensureAttachments() *platform.Attachment {
	if m.Attachments == nil {
		m.Attachments = []*platform.Attachment{&platform.Attachment{}}
	}
	a := m.Attachments[0]
	if a.Fallback == "" {
		a.Fallback = "This message is not supported by your client."
	}
	return a
}

func (m *MattermostMessage) ensureFields() *platform.Attachment {
	a := m.ensureAttachments()
	if a.Fields == nil {
		a.Fields = []*platform.AttachmentField{}
	}
	return a
}
