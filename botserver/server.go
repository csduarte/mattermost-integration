package botserver

import (
	"fmt"
	"net/http"
)

// Server represents the entire server process
type Server struct {
	*Config
	Triggers map[string]Trigger
	Store    *store
	Mux      *Mux
}

// NewServer takes config path and returns Server
func NewServer(configPath string) (*Server, error) {
	server := Server{}
	server.Mux = NewMux()
	config, err := parseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Config Parse Error - %v", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	server.Config = &config
	if err = server.LoadTriggers(); err != nil {
		return nil, err
	}

	return &server, nil
}

// LoadTriggers takes config triggers and loads them as actual webhooks
func (s *Server) LoadTriggers() error {
	if s.Config.Triggers == nil {
		return fmt.Errorf("Can't load triggers without a config")
	}
	if len(s.Config.Triggers) == 0 {
		return fmt.Errorf("Can't load triggers without triggers")
	}
	for k, v := range s.Config.Triggers {
		if s.Triggers == nil {
			s.Triggers = map[string]Trigger{}
		}
		s.Triggers[k] = NewTriggerFromConfig(k, v)
	}
	return nil
}

// Start Server will start a server (TLS Server if both the
// TLSCert and TLSKey config values are set) that will listen to all
// routes and redistribute the requests to the proper Integration
// handlers
func (s *Server) Start() error {
	listenAddr := fmt.Sprintf("%s:%s", s.Config.Host.BindAddress, s.Config.Host.BindPort)

	http.HandleFunc("/", s.primaryHandler)

	for k, v := range s.Mux.m {
		fmt.Printf("%v: %v\n", k, v)
	}

	var err error
	if s.Config.Host.TLSCert != "" && s.Config.Host.TLSKey != "" {
		fmt.Println("Starting HTTPS server on", listenAddr)
		err = http.ListenAndServeTLS(listenAddr,
			s.Config.Host.TLSCert,
			s.Config.Host.TLSKey,
			nil)
	} else {
		fmt.Println("Starting Unsecure HTTP server on", listenAddr)
		err = http.ListenAndServe(listenAddr, nil)
	}

	if err != nil {
		return err
	}

	return nil
}

// HandleAllChatter adds handler for all chat services that match incoming route
// and command pattern
func (s *Server) HandleAllChatter(pattern string, fn Handler) {
	for _, t := range s.Triggers {
		if t.isChatter() {
			s.Mux.add(t, pattern, fn)
		}
	}
}

// HandleSomeChatter adds handler for an array of chat services that match
// incoming route and command pattern
func (s *Server) HandleSomeChatter(names []string, pattern string, fn Handler) {
	for _, n := range names {
		t, ok := s.Triggers[n]
		if !ok {
			panic(fmt.Sprintf("HandleSome failed to find trigger %q", n))
		}
		if !t.isChatter() {
			panic(fmt.Sprintf("HandleSome tried to set pattern on non-chatter %q", n))
		}
		s.Mux.add(t, pattern, fn)
	}
}

// HandleChatter adds handler for matching chat service that match incoming
// route and command pattern
func (s *Server) HandleChatter(name, pattern string, fn Handler) {
	t, ok := s.Triggers[name]
	if !ok {
		panic(fmt.Sprintf("HandleOne failed to find trigger %q", name))
	}
	if !t.isChatter() {
		panic(fmt.Sprintf("HandleOne tried to set pattern on non-chatter %q", name))
	}
	s.Mux.add(t, pattern, fn)
}
