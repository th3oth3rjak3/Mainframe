package api

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/th3oth3rjak3/mainframe/internal/docs"
	"github.com/th3oth3rjak3/mainframe/internal/domain"
	"github.com/th3oth3rjak3/mainframe/internal/handler"
	mw "github.com/th3oth3rjak3/mainframe/internal/middleware"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

// Server holds the dependencies for the HTTP server.
type Server struct {
	router    *fiber.App
	container *ServiceContainer
	hmacKey   string
}

// NewServer creates a new Server instance and configures its routes.
func NewServer(container *ServiceContainer, hmacKey string, webAssets embed.FS) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler:          customErrorHandler,
		Immutable:             true, // Context safety!
		ReadTimeout:           60 * time.Second,
		WriteTimeout:          60 * time.Second,
		IdleTimeout:           60 * time.Second,
		DisableStartupMessage: false, // Keep that nice Fiber banner
	})

	// Create the server instance
	s := &Server{
		container: container,
		router:    app,
		hmacKey:   hmacKey,
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080, http://127.0.0.1:8080",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Attach middleware
	s.router.Use(logger.New())
	s.router.Use(recover.New())

	// Register routes
	s.registerRoutes()

	s.router.Use(filesystem.New(filesystem.Config{
		Root:         http.FS(webAssets), // embed FS
		PathPrefix:   "web",              // match everything else
		Index:        "index.html",
		NotFoundFile: "web/index.html", // SPA fallback
		Browse:       false,
	}))

	return s
}

// Start runs the HTTP server on the given address.
func (s *Server) Start(addr string) error {
	return s.router.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.router.ShutdownWithContext(ctx)
}

// registerRoutes sets up all the HTTP routes for the application.
func (s *Server) registerRoutes() {
	// Simple health check route
	s.router.Get("/health", handler.HandleHealthCheck)

	// Documentation routes
	s.router.Get("/openapi.json", handler.HandleOpenAPI)
	s.router.Get("/docs", handler.HandleDocs)

	// API routes
	apiGroup := s.router.Group("/api")
	authGroup := apiGroup.Group("/auth")

	// auth routes
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return handler.HandleLogin(c, s.container.AuthenticationService, s.container.CookieService)
	})

	// PROTECTED ROUTES
	authMiddleware := mw.NewAuthMiddleware(
		s.container.SessionRepository,
		s.container.UserRepository,
		s.container.CookieService,
		s.hmacKey,
	)

	protectedGroup := apiGroup.Group("", authMiddleware.SessionAuth)

	adminRoleRequired := mw.RequireRole(domain.Administrator)

	// logout route
	authGroup.Post("/logout",
		authMiddleware.SessionAuth,
		func(c *fiber.Ctx) error {
			return handler.HandleLogout(c, s.container.AuthenticationService, s.container.CookieService)
		})

	// Users Group
	usersGroup := protectedGroup.Group("/users", adminRoleRequired)
	usersGroup.Get("", func(c *fiber.Ctx) error {
		return handler.HandleListUsers(c, s.container.UserService)
	})
	usersGroup.Get("/:id", func(c *fiber.Ctx) error {
		return handler.HandleGetUserByID(c, s.container.UserService)
	})
	usersGroup.Post("", func(c *fiber.Ctx) error {
		return handler.HandleCreateUser(c, s.container.UserService)
	})
	usersGroup.Put("/:id", func(c *fiber.Ctx) error {
		return handler.HandleUpdateUser(c, s.container.UserService)
	})

	// Roles group
	rolesGroup := protectedGroup.Group("/roles", adminRoleRequired)
	rolesGroup.Get("", func(c *fiber.Ctx) error {
		return handler.HandleListRoles(c, s.container.RoleService)
	})
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Check for fiber errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	// check validation errorss
	if _, ok := err.(validation.Errors); ok {
		code = fiber.StatusBadRequest
	}

	// Check domain error types
	switch {
	case errors.Is(err, shared.ErrNotFound):
		code = fiber.StatusNotFound
	case errors.Is(err, shared.ErrUsernameTaken):
		code = fiber.StatusConflict
	case errors.Is(err, shared.ErrForbidden):
		code = fiber.StatusForbidden
	case errors.Is(err, shared.ErrUnauthorized), errors.Is(err, shared.ErrInvalidCredentials):
		code = fiber.StatusUnauthorized
	case errors.Is(err, shared.ErrBadRequest):
		code = fiber.StatusBadRequest
	}

	// Return JSON response
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
