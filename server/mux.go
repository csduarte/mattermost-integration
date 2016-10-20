package server

import "sync"

// Mux takes the singles from the server and maps to muxEntry
type Mux struct {
	mu sync.RWMutex
	m  map[string][]*muxEntry
}

type muxEntry struct {
	h            Handler
	pattern      string
	integrations *Integration
}

// Handler will perform action based to an incoming webhook context
type Handler func(context *Context)

// NewMux Generates a new blank Mux
func NewMux() *Mux {
	d := Mux{sync.RWMutex{}, make(map[string][]*muxEntry)}
	return &d
}

func (d *Mux) add(in *Integration, pattern string, h Handler) {
	d.mu.Lock()
	n := in.Name
	if d.m[n] == nil {
		d.m[n] = []*muxEntry{}
	}
	d.m[n] = append(d.m[n], &muxEntry{h, pattern, in})
	d.mu.Unlock()
}
