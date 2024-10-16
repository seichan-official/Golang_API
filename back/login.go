package main

import (
    "net/http"

    "github.com/labstack/echo"

    "github.com/gericass/gotify"
)

var Auth gotify.OAuth
var Token gotify.Gotify


func Handler(c echo.Context) error {
    url := Auth.AuthURL() // Get the Redirect URL for authenticate
    return c.Redirect(301, url)
}

// CallbackHandler : Controller for https://localhost:3000/callback/
func CallbackHandler(c echo.Context) error {

    t, err := Auth.GetToken(c.Request()) // Get the token for using Spotify API
    if err != nil {
        return err
    }
    Token = t

    return c.String(http.StatusOK, "Authentication success")
}

// RefreshHandler : Controller for https://localhost:3000/refresh/
func RefreshHandler(c echo.Context) error {

    err := Token.Refresh() // Refreshing token for using Spotify API
    if err != nil {
        return err
    }

    return c.String(http.StatusOK, "AccessToken Refreshed")
}
