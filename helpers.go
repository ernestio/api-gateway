package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	ErrUnauthorized   = echo.NewHTTPError(http.StatusForbidden, "")
	ErrNotFound       = echo.NewHTTPError(http.StatusNotFound, "")
	ErrBadReqBody     = echo.NewHTTPError(http.StatusBadRequest, "something")
	ErrGatewayTimeout = echo.NewHTTPError(http.StatusGatewayTimeout, "")
	ErrInternal       = echo.NewHTTPError(http.StatusInternalServerError, "")
)

// Get the authenticated user from the JWT Token
func authenticatedUser(c echo.Context) User {
	var u User

	user := c.Get("user").(*jwt.Token)
	u.Username = user.Claims["username"].(string)
	u.Admin = user.Claims["admin"].(bool)

	return u
}
