package proxy

import (
	"fmt"
	"log"
	"net/http"
)

// Server is the sandgrouse proxy server.
type Server struct {
	ListenAddr string
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "sandgrouse proxy running")
	})

	log.Printf("sandgrouse proxy listening on %s", s.ListenAddr)
	return http.ListenAndServe(s.ListenAddr, mux)
}
