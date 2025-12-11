package main

import (
	"log"
	"os"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/th3oth3rjak3/mainframe/docs"
	"github.com/th3oth3rjak3/mainframe/internal/handler"
	"github.com/th3oth3rjak3/mainframe/internal/repository"
)

func initDB() (*sqlx.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/mainframe.db"
	}
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable foreign keys for SQLite
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// @title           Mainframe API
// @version         1.0
// @description     User authentication and session management
// @host            localhost:8080
// @BasePath        /
func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	e := echo.New()

	// add basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(60 * time.Second))

	// simple health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "ok",
		})
	})

	// Serve the OpenAPI JSON
	e.GET("/openapi.json", func(c echo.Context) error {
		return c.File("./docs/swagger.json")
	})

	e.GET("/docs", func(c echo.Context) error {
		html, err := scalargo.NewV2(
			scalargo.WithSpecURL("http://localhost:8080/openapi.json"),
			scalargo.WithTheme(scalargo.ThemeDeepSpace),
		)
		if err != nil {
			return c.String(500, err.Error())
		}
		return c.HTML(200, html)
	})

	e.POST("/api/auth/login", func(c echo.Context) error {
		return handler.HandleLogin(c, userRepo)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
