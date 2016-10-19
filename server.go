package integrationserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/csduarte/integrationserver/platform"
)

// IntegrationServer represents the entire server process
type IntegrationServer struct {
	Store  *integrationStore
	Client *platform.Client
	Config Config
}

// Handler will perform action based to an incoming webhook context
type Handler func(context Context)

// NewIntegrationServer takes config path and returns IntegrationServer
func NewIntegrationServer(configPath string) (*IntegrationServer, error) {
	IntegrationServer := IntegrationServer{}

	config, err := parseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Config Parse Error - %v", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	IntegrationServer.Config = config
	IntegrationServer.Store, err = restoreIntegrationStore(config.StorageLocation)
	if err != nil {
		return nil, err
	}

	err = IntegrationServer.Store.matchIntegrations(config)
	if err != nil {
		return nil, err
	}

	return &IntegrationServer, nil
}

// parseConfig unmarshals json data for Config
func parseConfig(path string) (Config, error) {
	config := Config{}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return config, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

// Start IntegrationServer Server listening for webhooks
func (i *IntegrationServer) Start() error {
	http.HandleFunc("/", handler)
	listenAddr := fmt.Sprintf("%s:%s", i.Config.BindAddress, i.Config.BindPort)
	fmt.Println("Opening server on", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	// http.ListenAndServeTLS("127.0.0.1", "cert.pem", "key.pem", nil)
	if err != nil {
		return err
	}
	return nil
}

// Handle adds handler for all chat services
func (i *IntegrationServer) Handle(rs string, fn Handler) {

}

// HandleOne adds handler for matching chat services
func (i *IntegrationServer) HandleOne(n, rs string, fn Handler) {

}
