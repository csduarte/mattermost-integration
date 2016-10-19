package integrationserver

import "fmt"

// Config holds the process setup and integrations
type Config struct {
	Host            string
	BindAddress     string
	BindPort        string
	TLSCert         string
	TLSKey          string
	StorageLocation string
	Integrations    []integrationConfig
}

// integrationConfig holds information on a specfic integration
type integrationConfig struct {
	Name           string
	Server         string
	CreateHooks    bool
	TeamName       string
	Username       string
	Password       string
	ChannelName    string
	IncomingRoutes []string
	TriggerWords   []string
	TriggerWhen    string
}

func (i *integrationConfig) validate() error {
	if i.Server == "" {
		return fmt.Errorf("Integration(%q) mising host", i.Name)
	}
	if i.Username == "" {
		return fmt.Errorf("Integration(%q) mising username", i.Name)
	}
	if i.Password == "" {
		return fmt.Errorf("Integration(%q) mising password", i.Name)
	}
	if len(i.IncomingRoutes) < 1 {
		return fmt.Errorf("Integration(%q) has no IncomingRoutes", i.Name)
	}
	if i.TriggerWhen == "" {
		return fmt.Errorf("Integration(%q) mising TriggerWhen", i.Name)
	}
	if len(i.TriggerWords) < 1 {
		return fmt.Errorf("Integration(%q) missing TriggerWords", i.Name)
	}
	return nil
}

func (c *Config) validate() error {
	if c.Host == "" {
		return fmt.Errorf("IntegrationServer missing Host")
	}

	if c.BindAddress == "" {
		return fmt.Errorf("IntegrationServer missing BindAddress")
	}

	if c.BindPort == "" {
		return fmt.Errorf("IntegrationServer missing BindPort")
	}

	if len(c.Integrations) < 1 {
		return fmt.Errorf("No integrations listed in Config")
	}

	for _, integration := range c.Integrations {
		err := integration.validate()
		if err != nil {
			return err
		}
	}

	return nil
}
