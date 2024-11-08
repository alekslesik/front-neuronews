package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/alekslesik/neuro-news/internal/pkg/router"
	"github.com/alekslesik/neuro-news/pkg/config"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type Server struct {
	*http.Server
}

func New(c *config.Config, l *logger.Logger, r *router.Router) (*Server, error) {
	const op = "server.New()"

	// Get root certificate from system storage
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		l.Err(err).Msgf("%s > get root certificate", op)
		return nil, err
	}

	// Set up server
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", c.App.Host, c.App.Port),
		Handler: r.Route(),
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			RootCAs:            rootCAs,
			InsecureSkipVerify: false,
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{s}, nil
}
