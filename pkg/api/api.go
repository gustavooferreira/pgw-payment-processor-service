package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/api/middleware"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core/log"
)

// Server is the webserver environment, which holds all its dependencies.
type Server struct {
	Logger     log.Logger
	Repo       core.CreditCardChecker
	Authoriser core.Authoriser

	Router     *gin.Engine
	HTTPServer http.Server
}

// NewServer creates a new server.
func NewServer(addr string, port int, devMode bool, logger log.Logger, repo core.CreditCardChecker, authoriser core.Authoriser) *Server {
	s := &Server{Logger: logger, Repo: repo, Authoriser: authoriser}

	if !devMode {
		gin.SetMode(gin.ReleaseMode)
	}

	s.Router = gin.New()

	s.Router.Use(
		middleware.GinReqLogger(logger, time.RFC3339, "request served", "http-router-mux"),
	)
	if !devMode {
		s.Router.Use(gin.Recovery())
	}

	// Create http.Server
	s.HTTPServer = http.Server{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.setupRoutes(devMode)

	return s
}

// setupRoutes creates routes for all handlers
func (s *Server) setupRoutes(devMode bool) {
	s.Router.NoRoute(NoRoute)
	v1 := s.Router.Group("/api/v1")
	v1.GET("/healthcheck", s.Healthcheck)

	v1.POST("/authorise", s.AuthoriseTransaction)
	v1.POST("/capture", s.CaptureTransaction)
	v1.POST("/void", s.VoidTransaction)
	v1.POST("/refund", s.RefundTransaction)

	// Profiler
	// URL: https://<IP>:<PORT>/debug/pprof/
	if devMode {
		s.Logger.Info("activating pprof (devmode on)", log.Field("type", "debug"))
		pprof.Register(s.Router)
	}
}

// ListenAndServe listens and serves incoming requests.
func (s *Server) ListenAndServe() error {
	if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// ShutDown gracefully shuts down server.
func (s *Server) ShutDown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
