package integrationserver

import (
	"fmt"

	"github.com/csduarte/integrationserver/platform"
)

// Integration holds incoming and outgoing webhook from mattermost
type Integration struct {
	Name       string                   `json:"name"`
	Host       string                   `json:"host"`
	ChatServer string                   `json:"server"`
	ToMM       platform.IncomingWebhook `json:"toMM"`
	FromMM     platform.OutgoingWebhook `json:"fromMM"`
	Config     integrationConfig        `json:"config"`
}

// NewIntegrationFromConfig takes a integrationConfig and creates mm webhooks
func NewIntegrationFromConfig(s Config, ic integrationConfig) Integration {
	i := Integration{}
	i.Name = ic.Name
	i.Host = s.Host
	i.ChatServer = ic.Server
	i.Config = ic
	return i
}

func (i *Integration) initialize(c platform.Client) error {
	err := i.createFromHook(c)
	if err != nil {
		return err
	}
	err = i.createToHook(c)
	if err != nil {
		return err
	}
	return nil
}

func (i *Integration) createToHook(c platform.Client) error {
	fmt.Println("Creating To Hook")
	wh := platform.IncomingWebhook{}
	channelID, err := c.FindChannelIDByName(
		i.Config.ToMattermost.ChannelName,
		i.Config.TeamName,
	)
	if err != nil {
		return err
	}
	wh.ChannelID = channelID
	res, rErr := c.CreateIncomingWebhook(i.Config.TeamName, &wh)
	if rErr != nil {
		return fmt.Errorf("Failed to createToHook for team %q - %v", i.Config.TeamName, rErr.Message)
	}
	i.ToMM = *res.Data.(*platform.IncomingWebhook)
	return nil
}

func (i *Integration) createFromHook(c platform.Client) error {
	fmt.Println("Creating From Hooks")
	wh := platform.OutgoingWebhook{}

	wh.DisplayName = fmt.Sprintf("(IntegrationServer) %v", i.Name)
	wh.TriggerWords = i.Config.FromMattermost.TriggerWords
	wh.CallbackURLs = i.createCallbacks()

	res, err := c.CreateOutgoingWebhook(i.Config.TeamName, &wh)
	if err != nil {
		return fmt.Errorf("Failed to createFromHook for team %q - %v", i.Config.TeamName, err.Message)
	}
	i.FromMM = *res.Data.(*platform.OutgoingWebhook)
	return nil
}

func (i *Integration) createCallbacks() []string {
	callbacks := make([]string, len(i.Config.FromMattermost.IncomingRoutes))
	for idx, v := range i.Config.FromMattermost.IncomingRoutes {
		callbacks[idx] = i.Host + v
	}
	fmt.Println("Callback Count:", len(callbacks))
	return callbacks
}
