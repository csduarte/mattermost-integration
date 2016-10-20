package server

import "fmt"

// Config holds the process setup and integrations
type Config struct {
	Host         string
	BindAddress  string
	BindPort     string
	TLSCert      string
	TLSKey       string
	Integrations []*integrationConfig
}

// integrationConfig holds information on a specfic integration
type integrationConfig struct {
	Name           string
	Server         string
	TeamName       string
	Username       string
	Password       string
	DeleteRemote   bool
	ToMattermost   *toMMConfig
	FromMattermost *fromMMConfig
}

type fromMMConfig struct {
	IncomingRoutes []string
	TriggerWords   []string
	TriggerExact   bool
}

type toMMConfig struct {
	ChannelName string
}

func (i *integrationConfig) validate() error {
	if i.Name == "" {
		return fmt.Errorf("Integration(%q) mising name", i.Name)
	}
	if i.Server == "" {
		return fmt.Errorf("Integration(%q) mising host", i.Name)
	}
	if i.TeamName == "" {
		return fmt.Errorf("Integration(%q) mising team name", i.Name)
	}
	if i.Username == "" {
		return fmt.Errorf("Integration(%q) mising username", i.Name)
	}
	if i.Password == "" {
		return fmt.Errorf("Integration(%q) mising password", i.Name)
	}
	if i.ToMattermost != nil {
		if err := i.ToMattermost.validate(); err != nil {
			return fmt.Errorf("Integration(%q) - ", err.Error())
		}
	}
	if i.FromMattermost != nil {
		if err := i.FromMattermost.validate(); err != nil {
			return fmt.Errorf("Integration(%q) - ", err.Error())
		}
	}
	return nil
}

func (f *fromMMConfig) validate() error {
	if f.IncomingRoutes == nil || len(f.IncomingRoutes) < 1 {
		return fmt.Errorf("Missing Incoming Routes")
	}
	if f.TriggerWords == nil || len(f.TriggerWords) < 1 {
		return fmt.Errorf("Missing Trigger Words")
	}
	for _, ir := range f.IncomingRoutes {
		if ir == "" {
			return fmt.Errorf("No empty routes")
		}
	}
	for _, ir := range f.TriggerWords {
		if ir == "" {
			return fmt.Errorf("No empty trigger words")
		}
	}
	return nil
}

func (t *toMMConfig) validate() error {
	if t.ChannelName == "" {
		return fmt.Errorf("Empty Channel name")
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
