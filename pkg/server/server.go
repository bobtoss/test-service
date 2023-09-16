package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	http *http.Server
}

type Configuration func(r *Server) error

func New(configs ...Configuration) (r *Server, err error) {
	r = &Server{}

	for _, cfg := range configs {
		if err = cfg(r); err != nil {
			return
		}
	}
	return
}

func (s *Server) Run(logger *zap.Logger) (err error) {
	if s.http != nil {
		go func() {
			if err = s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("ERR_INIT_REST", zap.Error(err))
				return
			}
		}()
		logger.Info("http server started on http://localhost" + s.http.Addr)
	}
	return
}

func (s *Server) Stop(ctx context.Context) (err error) {
	if s.http != nil {
		if err = s.http.Shutdown(ctx); err != nil {
			return
		}
	}

	return
}

func WithHTTPServer(handler http.Handler, port string) Configuration {
	return func(s *Server) (err error) {
		s.http = &http.Server{
			Handler: handler,
			Addr:    ":" + port,
		}
		return
	}
}
