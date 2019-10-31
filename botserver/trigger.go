package botserver

const (
	// MattermostTriggerType holds string name for trigger
	MattermostTriggerType TriggerType = "Mattermost"
	// ExternalTriggerType holds string name for trigger
	ExternalTriggerType TriggerType = "External"
)

// TriggerType well tell system how account for speicfic settings
type TriggerType string

// Trigger represents the various possible triggers
type Trigger interface {
	Name() string
	isChatter() bool
}

// MattermostTrigger represents a mattermost outgoingWebhook
type MattermostTrigger struct {
	KeyName string
	Type    TriggerType
	Config  *TriggerConfig
}

// ExternalTrigger represents an external webhook from some other web services
type ExternalTrigger struct {
	Type TriggerType
}

// NewTriggerFromConfig will create a new trigger based on type
func NewTriggerFromConfig(key string, tc *TriggerConfig) Trigger {
	switch tc.Type {
	case string(MattermostTriggerType):
		return NewMattermostTriggerFromConfig(key, tc)
	case string(ExternalTriggerType):
		return NewExternalTriggerFromConfig(key, tc)
	default:
		return NewMattermostTriggerFromConfig(key, tc) // going to return something...
	}
}

// NewMattermostTriggerFromConfig will create a new trigger for mattermost
func NewMattermostTriggerFromConfig(key string, tc *TriggerConfig) *MattermostTrigger {
	t := &MattermostTrigger{}
	t.Type = MattermostTriggerType
	t.Config = tc
	t.KeyName = key
	return t
}

// NewExternalTriggerFromConfig will create a new trigger for mattermost
func NewExternalTriggerFromConfig(key string, tc *TriggerConfig) *MattermostTrigger {
	t := &MattermostTrigger{}
	t.KeyName = key
	return t
}

// Name returns MattermostTrigger name
func (t *MattermostTrigger) Name() string {
	return t.KeyName
}

// isChatter returns true if trigger is of a chat server type
func (t *MattermostTrigger) isChatter() bool {
	return t.Type == MattermostTriggerType
}
