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
	Config
	Client *platform.Client
	Store  *integrationStore
	Demux  *Demux
}

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
	IntegrationServer.Store, err = restoreIntegrationStore(StorageLocation)
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

func (i *IntegrationServer) muxHandler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	defer req.Body.Close()
	fmt.Println(path)
	// match path
	// match integrations
	// match trigger words
	// build context
	// pass to matching handler
}

// Start IntegrationServer will start a server (TLS Server if both the
// TLSCert and TLSKey config values are set) that will listen to all
// routes and redistribute the requests to the proper Integration
// handlers
func (i *IntegrationServer) Start() error {
	listenAddr := fmt.Sprintf("%s:%s", i.Config.BindAddress, i.Config.BindPort)

	http.HandleFunc("/", i.muxHandler)

	var err error
	if i.TLSCert != "" && i.TLSKey != "" {
		fmt.Println("Starting HTTPS server on", listenAddr)
		err = http.ListenAndServeTLS(listenAddr, i.TLSCert, i.TLSKey, nil)
	} else {
		fmt.Println("Starting HTTP server on", listenAddr)
		err = http.ListenAndServe(listenAddr, nil)
	}

	if err != nil {
		return err
	}

	return nil
}

// HandleAll adds handler for all chat services that match incoming route
// and message pattern
func (i *IntegrationServer) HandleAll(pattern string, fn Handler) {
	integrations := i.Store.Integrations
	for _, in := range integrations {
		i.Demux.add(in, pattern, fn)
	}
}

// HandleSome adds handler for an array of chat services that match
// incoming route and message pattern
func (i *IntegrationServer) HandleSome(names []string, pattern string, fn Handler) {
	for _, n := range names {
		i.Demux.add(i.Store.findByName(n), pattern, fn)
	}
}

// HandleOne adds handler for matching chat service that match incoming
// route and message pattern
func (i *IntegrationServer) HandleOne(name, pattern string, fn Handler) {
	integration := i.Store.findByName(name)
	i.Demux.add(integration, pattern, fn)
}
