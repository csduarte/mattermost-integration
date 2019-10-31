package botserver

import "net/http"

func (s *Server) primaryHandler(w http.ResponseWriter, r *http.Request) {
	// 	defer r.Body.Close()
	// 	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte("Bad Content Type"))
	// 		fmt.Printf("Bad Content Type\n")
	// 		return
	// 	}
	//
	// 	c := NewContext(w, r)
	// 	if c.payload == nil {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("Unexpected Payload"))
	// 		fmt.Printf("Unexpected Payload\n")
	// 		return
	// 	}
	//
	// 	in := i.Store.findByToken(c.payload.Token)
	// 	if in == nil {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("Bad Token"))
	// 		fmt.Printf("Bad Token\n")
	// 		return
	// 	}
	// 	c.addIntegration(in)
	//
	// 	routeMatches := false
	// 	for _, ir := range in.Config.FromMattermost.IncomingRoutes {
	// 		if ir == r.URL.Path {
	// 			routeMatches = true
	// 		}
	// 	}
	// 	if !routeMatches {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("No Route Matching: " + r.URL.Path))
	// 		fmt.Printf("No Route Matching: %v", r.URL.Path)
	// 		return
	// 	}
	//
	// 	if !in.FromMM.HasTriggerWord(c.payload.TriggerWord) {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("No Matching Trigger"))
	// 		fmt.Printf("No Matching Trigger")
	// 		return
	// 	}
	//
	// 	commands := strings.Fields(c.payload.Text)
	// 	if len(commands) < 1 {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write([]byte("No text content"))
	// 		fmt.Printf("No text content")
	// 		return
	// 	}
	// 	var command string
	// 	if len(commands) < 2 {
	// 		command = ""
	// 	} else if commands[1] != c.payload.TriggerWord {
	// 		command = commands[1]
	// 	} else {
	// 		command = commands[1]
	// 	}
	//
	// 	entries, ok := i.Mux.m[in.Name]
	// 	if !ok {
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte("No Commands Registered for Integration"))
	// 		fmt.Printf("No Commands Registered for Integration")
	// 		return
	// 	}
	//
	// 	for _, ent := range entries {
	//
	// 		if m := ent.pattern.Match([]byte(command)); !m {
	// 			continue
	// 		}
	//
	// 		ent.h(c)
	// 		if c.response != nil {
	// 			w.WriteHeader(http.StatusOK)
	// 			w.Write([]byte(c.response.ToJSON()))
	// 		} else {
	// 			w.WriteHeader(http.StatusOK)
	// 			w.Write([]byte("Ok"))
	// 		}
	//
	// 		return
	// 	}
	//
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("No Matching Command"))
	// 	fmt.Printf("No Matching Command")
	//
}
