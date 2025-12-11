package handler

import (
	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/labstack/echo/v4"
	_ "github.com/th3oth3rjak3/mainframe/docs"
)

func HandleOpenAPI(c echo.Context) error {
	return c.File("./docs/swagger.json")
}

func HandleDocs(c echo.Context) error {
	html, err := scalargo.NewV2(
		scalargo.WithSpecURL("http://localhost:8080/openapi.json"),
		scalargo.WithTheme(scalargo.ThemeDeepSpace),
	)
	if err != nil {
		return c.String(500, err.Error())
	}
	return c.HTML(200, html)
}
