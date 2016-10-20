package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// IntegrationServer represents the entire server process
type IntegrationServer struct {
	Config
	Store *integrationStore
	Mux   *Mux
	// Client *platform.Client
}

// NewIntegrationServer takes config path and returns IntegrationServer
func NewIntegrationServer(configPath string) (*IntegrationServer, error) {
	server := IntegrationServer{}
	server.Mux = NewMux()
	config, err := parseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Config Parse Error - %v", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	server.Config = config
	server.Store, err = restoreIntegrationStore(StorageLocation)
	if err != nil {
		return nil, err
	}

	err = server.Store.matchIntegrations(config)
	if err != nil {
		return nil, err
	}

	return &server, nil
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

func (i *IntegrationServer) primaryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Content Type"))
		fmt.Printf("Bad Content Type\n")
		return
	}

	c := NewContext(w, r)
	if c.payload == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unexpected Payload"))
		fmt.Printf("Unexpected Payload\n")
		return
	}

	in := i.Store.findByToken(c.payload.Token)
	if in == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Bad Token"))
		fmt.Printf("Bad Token\n")
		return
	}
	c.addIntegration(in)

	routeMatches := false
	for _, ir := range in.Config.FromMattermost.IncomingRoutes {
		if ir == r.URL.Path {
			routeMatches = true
		}
	}
	if !routeMatches {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Route Matching: " + r.URL.Path))
		fmt.Printf("No Route Matching: %v", r.URL.Path)
		return
	}

	if !in.FromMM.HasTriggerWord(c.payload.TriggerWord) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Matching Trigger"))
		fmt.Printf("No Matching Trigger")
		return
	}

	commands := strings.Fields(c.payload.Text)
	if len(commands) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No text content"))
		fmt.Printf("No text content")
		return
	}
	var command string
	if len(commands) < 2 {
		command = ""
	} else if commands[1] != c.payload.TriggerWord {
		command = commands[1]
	} else {
		command = commands[1]
	}

	entries, ok := i.Mux.m[in.Name]
	if !ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No Commands Registered for Integration"))
		fmt.Printf("No Commands Registered for Integration")
		return
	}

	for _, ent := range entries {

		m, err := regexp.Match(ent.pattern, []byte(command))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad Pattern Registered"))
			fmt.Printf("Bad Pattern Registered - %v", ent.pattern)
			return
		}

		if !m {
			continue
		}

		ent.h(c)
		if c.response != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(c.response.ToJSON()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ok"))
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("No Matching Command"))
	fmt.Printf("No Matching Command")
}

// Start IntegrationServer will start a server (TLS Server if both the
// TLSCert and TLSKey config values are set) that will listen to all
// routes and redistribute the requests to the proper Integration
// handlers
func (i *IntegrationServer) Start() error {
	listenAddr := fmt.Sprintf("%s:%s", i.Config.BindAddress, i.Config.BindPort)

	http.HandleFunc("/", i.primaryHandler)

	for k, v := range i.Mux.m {
		fmt.Printf("%v: %v\n", k, v)
	}

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
		i.Mux.add(in, pattern, fn)
	}
}

// HandleSome adds handler for an array of chat services that match
// incoming route and message pattern
func (i *IntegrationServer) HandleSome(names []string, pattern string, fn Handler) {
	for _, n := range names {
		in := i.Store.findByName(n)
		if in == nil {
			panic(fmt.Sprintf("HandleSome failed to find Integration %q", n))
		}
		i.Mux.add(in, pattern, fn)
	}
}

// HandleOne adds handler for matching chat service that match incoming
// route and message pattern
func (i *IntegrationServer) HandleOne(name, pattern string, fn Handler) {
	in := i.Store.findByName(name)
	if in == nil {
		panic(fmt.Sprintf("HandleOne failed to find Integration %q", name))
	}
	i.Mux.add(in, pattern, fn)
}
