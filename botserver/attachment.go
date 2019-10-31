package botserver

import "github.com/csduarte/mattermost-integration/platform"

// Attachment holds message attachments
type Attachment struct {
	pa *platform.Attachment
}

// NewAttachment generates attachment
func NewAttachment() *Attachment {
	a := Attachment{}
	a.pa = &platform.Attachment{}
	return &a
}

// SetAuthorIcon sets author icon
func (a *Attachment) SetAuthorIcon(url string) {
	a.pa.AuthorIcon = url
}

// SetAuthorLink sets author link
func (a *Attachment) SetAuthorLink(url string) {
	a.pa.AuthorLink = url
}

// SetAuthorName sets author name
func (a *Attachment) SetAuthorName(text string) {
	a.pa.AuthorName = text
}

// SetColor sets attachment color (ex: #333333)
func (a *Attachment) SetColor(cssHexColor string) {
	a.pa.Color = cssHexColor
}

// SetFallback sets attachment fallback text
func (a *Attachment) SetFallback(text string) {
	a.pa.Fallback = text
}

// SetPretext sets attachment Pretext
func (a *Attachment) SetPretext(text string) {
	a.pa.Pretext = text
}

// SetText sets attachment text
func (a *Attachment) SetText(text string) {
	a.pa.Text = text
}

// SetTitle sets attachment title
func (a *Attachment) SetTitle(text string) {
	a.pa.Title = text
}

// SetTitleLink sets attachment title link
func (a *Attachment) SetTitleLink(url string) {
	a.pa.TitleLink = url
}

// SetImageURL will add a image url attachment
func (a *Attachment) SetImageURL(url string) {
	a.pa.ImageURL = url
}

// AddField adds a field to Attachment
func (a *Attachment) AddField(title, text string, short bool) *platform.AttachmentField {
	f := platform.AttachmentField{}
	f.Short = short
	f.Title = title
	f.Value = text
	a.pa.Fields = append(a.pa.Fields, &f)
	return &f
}
