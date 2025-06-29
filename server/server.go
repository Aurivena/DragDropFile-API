package server

import (
	"DragDrop-Files/models"
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

const DEVELOPMENT = "development"

func (s *Server) Run(port string, routes http.Handler, certificates models.CertificatesConfig) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        routes,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
	}
	logrus.Info("server started successfully")
	return s.httpServer.ListenAndServeTLS(certificates.CertificatesPath, certificates.KeyPath)
}

func (s *Server) Stop(ctx context.Context, postgres *sqlx.DB) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Info("http listener shutdown successfully")
	}

	if err := postgres.Close(); err != nil {
		logrus.Error(err.Error())
	} else {
		logrus.Info("business database connection closed successfully")
	}

	logrus.Info("server shutdown process completed successfully")
}
