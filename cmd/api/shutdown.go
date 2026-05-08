package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server wraps the HTTP server with graceful shutdown capabilities
type Server struct {
	srv *http.Server
	log *application
}

// NewServer creates a new Server instance
func NewServer(srv *http.Server, log *application) *Server {
	return &Server{
		srv: srv,
		log: log,
	}
}

// ListenAndServeWithGracefulShutdown starts the server and handles graceful shutdown
func (s *Server) ListenAndServeWithGracefulShutdown() error {
	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to track server errors
	errChan := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		s.log.logger.Infow("Server starting", "addr", s.srv.Addr)
		errChan <- s.srv.ListenAndServe()
	}()

	// Wait for either signal or server error
	select {
	case sig := <-sigChan:
		s.log.logger.Infow("Received signal, shutting down gracefully", "signal", sig)
		return s.shutdown()
	case err := <-errChan:
		if err != http.ErrServerClosed {
			s.log.logger.Errorw("Server error", "error", err)
			return err
		}
		return nil
	}
}

// shutdown performs graceful shutdown with timeout
func (s *Server) shutdown() error {
	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.log.logger.Info("Graceful shutdown initiated")

	// Shutdown the server
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.logger.Errorw("Server shutdown error", "error", err)
		// Force close if graceful shutdown fails
		return s.srv.Close()
	}

	s.log.logger.Info("Server shutdown completed successfully")
	return nil
}
