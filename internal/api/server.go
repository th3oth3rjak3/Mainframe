package api

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/th3oth3rjak3/mainframe/internal/handler"
	mw "github.com/th3oth3rjak3/mainframe/internal/middleware"
)

// Server holds the dependencies for the HTTP server.
type Server struct {
	router    *echo.Echo
	container *ServiceContainer
}

// NewServer creates a new Server instance and configures its routes.
func NewServer(container *ServiceContainer) *Server {
	e := echo.New()
	// Create the server instance
	s := &Server{
		container: container,
		router:    e,
	}

	// Attach middleware
	s.router.Use(mw.ZerologRequestLogger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.ContextTimeout(60 * time.Second))

	// Register routes
	s.registerRoutes()

	return s
}

// Start runs the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	return s.router.Start(addr)
}

// registerRoutes sets up all the HTTP routes for the application.
func (s *Server) registerRoutes() {
	// Simple health check route
	s.router.GET("/health", handler.HandleHealthCheck)

	// Documentation routes
	s.router.GET("/openapi.json", handler.HandleOpenAPI)
	s.router.GET("/docs", handler.HandleDocs)

	// API routes
	apiGroup := s.router.Group("/api")
	authGroup := apiGroup.Group("/auth")

	// auth routes
	authGroup.POST("/login", func(c echo.Context) error {
		return handler.HandleLogin(c, s.container.AuthenticationService)
	})
}
