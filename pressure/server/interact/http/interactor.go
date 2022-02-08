package http

import (
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, handlers map[string]Handler) (*Server, []string, error) {
	dispatcher, err := toServerHandlers(handlers)
	if nil != err {
		return nil, nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/", dispatcher)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &Server{server: &server}, dispatcher.Patterns(), nil
}

func (s *Server) StartAsync(fErr func(error)) {
	go func() {
		err := s.server.ListenAndServe()
		if nil != err {
			fErr(err)
		}
	}()
}

func (s *Server) Close() error {
	return s.server.Close()
}
