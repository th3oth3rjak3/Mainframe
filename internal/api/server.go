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
	authMiddleware := mw.NewAuthMiddleware(
		s.container.SessionRepository,
		s.container.UserRepository,
		s.container.CookieService,
		s.hmacKey,
	)

	s.registerHealthCheckRoute()
	s.registerDocumentationRoutes()

	apiGroup := s.router.Group("/api")
	s.registerAuthenticationRoutes(apiGroup, authMiddleware)

	// Routes below here are all protected
	protectedGroup := apiGroup.Group("", authMiddleware.SessionAuth)
	s.registerUserRoutes(protectedGroup)
	s.registerRoleRoutes(protectedGroup)

}

// customErrorHandler is used in the fiber router to perform all error handling
// after the handler returns. This is used as a central location to handle
// error conversion to response status codes and return a consistent error
// type so it's easier to handle on the client.
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

// registerHealthCheckRoute associates the health check handler with the correct route.
func (s *Server) registerHealthCheckRoute() {
	s.router.Get("/health", handler.HandleHealthCheck)
}

// registerDocutnationRoutes registers the open api and scalar documentation endpoints
// with the server router.
func (s *Server) registerDocumentationRoutes() {
	s.router.Get("/openapi.json", handler.HandleOpenAPI)
	s.router.Get("/docs", handler.HandleDocs)
}

// registerAuthenticationRoutes registers all the routes associated with authentication
func (s *Server) registerAuthenticationRoutes(router fiber.Router, authMiddleware *mw.AuthMiddleware) {
	authGroup := router.Group("/auth")

	// Not protected on purpose to allow login
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		return handler.HandleLogin(c, s.container.AuthenticationService, s.container.CookieService)
	})

	authGroup.Post("/logout",
		authMiddleware.SessionAuth,
		func(c *fiber.Ctx) error {
			return handler.HandleLogout(c, s.container.AuthenticationService, s.container.CookieService)
		})

	authGroup.Get("/me", authMiddleware.SessionAuth, handler.HandleRefreshLoginDetails)
}

// registerUserRoutes registers all the routes associated with users.
// The router is expected to be protected by authentication middleware.
func (s *Server) registerUserRoutes(router fiber.Router) {
	adminRoleRequired := mw.RequireRole(domain.Administrator)
	usersGroup := router.Group("/users", adminRoleRequired)
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
	usersGroup.Delete("/:id", func(c *fiber.Ctx) error {
		return handler.HandleDeleteUser(c, s.container.UserService)
	})
}

// registerRoleRoutes registers all the routes associated with roles.
// The router is expectecd to be protected by authentication middleware.
func (s *Server) registerRoleRoutes(router fiber.Router) {
	adminRoleRequired := mw.RequireRole(domain.Administrator)
	rolesGroup := router.Group("/roles", adminRoleRequired)
	rolesGroup.Get("", func(c *fiber.Ctx) error {
		return handler.HandleListRoles(c, s.container.RoleService)
	})
}
