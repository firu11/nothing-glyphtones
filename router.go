package main

import (
	"gliphtones/database"
	"gliphtones/views"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, cmp templ.Component) error {
	//c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response())
}

func setupRouter(e *echo.Echo) {
	e.GET("/", index)
}

func index(c echo.Context) error {
	return Render(c, views.Index([]database.RingtoneModel{{Name: "sus", Id: 1}, {Name: "bak", Id: 2}}))
}
