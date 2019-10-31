package botserver

import "github.com/csduarte/mattermost-integration/platform"

// Message is a convience holder to represent Mattermost Incoming Webbhook
type Message struct {
	platform.IncomingWebhookRequest
	trigger *Trigger
}

// NewMessage initializes a new message
func NewMessage(channel string) *Message {
	m := Message{}
	m.ChannelName = channel
	return &m
}

// // Send message to associated integration
// func (m *Message) Send() error {
// 	if m.integration == nil {
// 		return fmt.Errorf("Could not send message - missing integration for username %q\n", m.Username)
// 	}
// 	if m.integration.ToMM != nil && m.integration.client != nil {
// 		in := m.integration
// 		data := m.IncomingWebhookRequest.ToJSON()
// 		fmt.Println(data)
// 		_, err := in.client.PostToWebhook(in.ToMM.ID, data)
// 		return fmt.Errorf(err.Message)
// 	}
// 	return fmt.Errorf("Could not send message - missing integration for %v\n", m.integration.Name)
// }
//
// // SetIconURL will simple add any string as the IconURL
// func (m *Message) SetIconURL(url string) {
// 	m.IconURL = url
// }
//
// // SetUsername sets username override, if allowed on server
// func (m *Message) SetUsername(name string) {
// 	m.Username = name
// }
//
// // SetMessage sets the text value
// func (m *Message) SetMessage(msg string) {
// 	m.Text = msg
// }
//
// // ensureAttachments make sure that the attachments has been set
// func (m *Message) ensureAttachments() *platform.Attachment {
// 	if m.Attachments == nil {
// 		m.Attachments = []*platform.Attachment{&platform.Attachment{}}
// 	}
// 	a := m.Attachments[0]
// 	if a.Fallback == "" {
// 		a.Fallback = "This message is not supported by your client."
// 	}
// 	return a
// }
//
// func (m *Message) ensureFields() *platform.Attachment {
// 	a := m.ensureAttachments()
// 	if a.Fields == nil {
// 		a.Fields = []*platform.AttachmentField{}
// 	}
// 	return a
// }
