package botserver

import "fmt"

// Config holds the process setup and integrations
type Config struct {
	Host     *HostConfig
	Servers  map[string]*ServerConfig
	Triggers map[string]*TriggerConfig
}

// HostConfig holds host server information
type HostConfig struct {
	URL         string
	BindAddress string
	BindPort    string
	TLSCert     string
	TLSKey      string
}

// ServerConfig holds information for server access
type ServerConfig struct {
	URL         string
	Username    string
	Password    string
	ChannelName string
}

// TriggerConfig holds information on a specfic Trigger
type TriggerConfig struct {
	Type string
	// MattermostTrigger
	ID             string
	ChannelName    string
	Description    string
	DisplayName    string
	IncomingRoutes []string
	Token          string
	TriggerExact   bool
	TriggerWords   []string
}

func (c *Config) validate() error {
	if c.Host == nil {
		return fmt.Errorf("Missing config `Host` key")
	}
	if c.Servers == nil {
		return fmt.Errorf("Missing config `Servers` key")
	}
	if c.Triggers == nil {
		return fmt.Errorf("Missing config `Triggers` key")
	}

	if err := c.Host.Valid(); err != nil {
		return err
	}

	if len(c.Servers) == 0 {
		return fmt.Errorf("Missing config `Servers` values")
	}

	if len(c.Triggers) == 0 {
		return fmt.Errorf("Missing config `triggers` values")
	}

	for k, v := range c.Servers {
		if err := v.Valid(); err != nil {
			return fmt.Errorf("Server %q - %v", k, err)
		}
	}

	for k, v := range c.Triggers {
		if err := v.Valid(); err != nil {
			return fmt.Errorf("Trigger %q - %v", k, err)
		}
	}

	return nil
}

// Valid will validate HostConfig
func (c *HostConfig) Valid() error {

	if c.URL == "" {
		return fmt.Errorf("Missing Host config `url`")
	}
	if c.BindAddress == "" {
		return fmt.Errorf("Missing Host config `BindAddress`")
	}
	if c.BindPort == "" {
		return fmt.Errorf("Missing Host config `BindPort`")
	}

	// TODO: c.TLSCert && c.TLSKey need to be valid

	return nil
}

// Valid will validate ServerConfig
func (c *ServerConfig) Valid() error {

	if c.URL == "" {
		return fmt.Errorf("Missing Host config `url`")
	}

	// TODO: Warng missing c.Username || c.Password
	// TODO: Warn missing c.Channelname

	return nil
}

// Valid will validate TriggerConfig
func (c *TriggerConfig) Valid() error {
	switch c.Type {
	case string(MattermostTriggerType):
		return c.validMattermostTrigger()
	default:
		return fmt.Errorf("Invalid Trigger Type %q", c.Type)
	}
}

func (c *TriggerConfig) validMattermostTrigger() error {

	if c.Token == "" {
		return fmt.Errorf("missing token")
	}

	if c.IncomingRoutes == nil || len(c.IncomingRoutes) < 1 {
		return fmt.Errorf("missing Incoming Routes")
	}
	if (c.TriggerWords == nil || len(c.TriggerWords) < 1) && c.ChannelName == "" {
		return fmt.Errorf("requres TriggerWords or ChannelName")
	}
	return nil
}
