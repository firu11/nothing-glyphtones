package main

import (
	"gliphtones/database"
	"gliphtones/templates/components"
	"gliphtones/templates/views"
	"log"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, cmp templ.Component) error {
	//c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response())
}

func setupRouter(e *echo.Echo) {
	e.RouteNotFound("/*", notFound)

	e.GET("/", index)
}

func index(c echo.Context) error {
	searchName := c.QueryParam("s")
	_ = c.QueryParam("s")

	// if it is a htmx request, render only one part
	if c.Request().Header.Get("HX-Request") == "true" {
		ringtones, err := database.GetRingtones(searchName, 0, 0)
		if err != nil {
			return err
		}
		return Render(c, components.ListOfRingtones(ringtones))
	}

	ringtones, err := database.GetRingtones(searchName, 0, 0)
	if err != nil {
		return err
	}

	var phones []database.PhoneModel
	if x, found := myCache.Get("phones"); found {
		log.Println("phones hit")
		phones = x.([]database.PhoneModel)
	} else {
		log.Println("phones miss")
		phones, err = database.GetPhones()
		if err != nil {
			return err
		}
		myCache.Set("phones", phones, 0)
	}

	var effects []database.EffectModel
	if x, found := myCache.Get("effects"); found {
		log.Println("effects hit")
		effects = x.([]database.EffectModel)
	} else {
		log.Println("effects miss")
		effects, err = database.GetEffects()
		if err != nil {
			return err
		}
		myCache.Set("effects", effects, 0)
	}

	return Render(c, views.Index(ringtones, phones, effects))
}

func notFound(c echo.Context) error {
	return Render(c, views.NotFound())
}

func google(c echo.Context) error {
	return nil
}
