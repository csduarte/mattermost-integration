package integrationserver

import "sync"

// Demux takes the singles from the server and maps to muxEntry
type Demux struct {
	mu sync.RWMutex
	m  map[string][]muxEntry
}

type muxEntry struct {
	h            Handler
	pattern      string
	integrations *Integration
}

// Handler will perform action based to an incoming webhook context
type Handler func(context Context)

func (d *Demux) add(in *Integration, pattern string, h Handler) {
	if d.m == nil {
		d.m = make(map[string][]muxEntry)
	}
	d.mu.Lock()
	for _, ir := range in.Config.FromMattermost.IncomingRoutes {
		d.m[ir] = append(d.m[ir], muxEntry{h, pattern, in})
	}
	d.mu.Unlock()
}
