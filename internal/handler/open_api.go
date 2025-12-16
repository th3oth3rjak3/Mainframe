package handler

import (
	_ "embed"
	"fmt"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/gofiber/fiber/v2"
)

func HandleOpenAPI(c *fiber.Ctx) error {
	return c.SendFile("./internal/docs/swagger.json")
}

func HandleDocs(c *fiber.Ctx) error {
	// Build full URL based on the request
	scheme := "http"
	if c.Protocol() == "https" {
		scheme = "https"
	}
	host := c.Hostname()
	fullURL := fmt.Sprintf("%s://%s/openapi.json", scheme, host)
	html, err := scalargo.NewV2(
		scalargo.WithSpecURL(fullURL),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
	)
	if err != nil {
		return err
	}
	c.Context().SetContentType("text/html")
	return c.SendString(html)
}
