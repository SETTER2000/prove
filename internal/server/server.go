// Package server implements HTTP server.
package server

import (
	"context"
	"fmt"
	"github.com/SETTER2000/prove/config"
	"github.com/SETTER2000/prove/internal/controller/grpc"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 50 * time.Second
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
	_defaultHTTPS           = false
	_defaultCertFile        = ""
	_defaultKeyFile         = ""
	_defaultGrpcPort        = ""
	//_defaultGrpcServer             = ""
)

// Server -.
type Server struct {
	server          *http.Server
	srv             *grpc.Server
	cfg             *config.Config
	notify          chan error
	certFile        string
	keyFile         string
	grpcPort        string
	shutdownTimeout time.Duration
	isHTTPS         bool
}

// New -.
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
		Addr:         _defaultAddr,
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
		//srv:             _defaultGrpcServer,
		shutdownTimeout: _defaultShutdownTimeout,
		isHTTPS:         _defaultHTTPS,
		certFile:        _defaultCertFile,
		keyFile:         _defaultKeyFile,
		grpcPort:        _defaultGrpcPort,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	switch s.isHTTPS {
	case true:
		go func() {
			s.srv.Logger.Info("Go backend serving on ", fmt.Sprintf("https://%s%s/ping", config.GetConfig().ServerDomain, s.server.Addr))
			s.notify <- s.server.ListenAndServeTLS(s.certFile, s.keyFile)
			close(s.notify)
		}()
	case false:
		go func() {
			s.notify <- s.server.ListenAndServe()
			close(s.notify)
		}()
	}

	go func() {
		s.srv.Logger.Info("Starting gRPC Server, PORT :", s.grpcPort)
		s.notify <- s.srv.ListenAndServer(s.grpcPort)
		close(s.notify)
	}()

}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
