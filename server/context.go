package server

import (
	"net/http"

	"github.com/csduarte/mattermost-integration/platform"
)

// Context hold request state
type Context struct {
	w        *http.ResponseWriter
	r        *http.Request
	i        *Integration
	response *platform.IncomingWebhookRequest
	payload  *platform.OutgoingWebhookPayload
}

// NewContext creates a integration context
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	c := Context{}
	c.w = &w
	c.r = r
	c.payload = platform.OutgoingWebhookPayloadFromForm(r.Body)
	return &c
}

// SeparateResponse will create a full featured ToMM webhook
func (c *Context) SeparateResponse() *MattermostMessage {
	m := NewMattermostMessage(c.payload.ChannelName)
	m.integration = c.i
	m.Username = c.i.Name
	return m
}

// SetIconURL to Response
func (c *Context) SetIconURL(url string) {
	c.ensureResponse()
	c.response.IconURL = url
}

// SetMessage to given context
func (c *Context) SetMessage(b string) {
	c.ensureResponse()
	c.response.Text = b
}

func (c *Context) addIntegration(i *Integration) {
	c.i = i
}

func (c *Context) ensureResponse() {
	if c.response == nil {
		c.response = &platform.IncomingWebhookRequest{}
		c.response.Username = "Bot!"
	}
}
