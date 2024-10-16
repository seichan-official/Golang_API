package main

import (
    "github.com/gericass/gotify"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

const clientID = ""
const clientSecret = ""
const callbackURI = "https://localhost:8000/callback/"

func main() {
    e := echo.New()
    e.Pre(middleware.HTTPSRedirect())
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    Auth = gotify.Set(clientID, clientSecret, callbackURI)

    e.GET("/", Handler)
    e.GET("/callback/", CallbackHandler)
    e.GET("/refresh/", RefreshHandler)

    // Require HTTPS
    e.Logger.Fatal(e.StartTLS(":8000", "cert.pem", "key.pem"))
}
