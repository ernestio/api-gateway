package helpers

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	// ErrUnauthorized : HTTP 403 error
	ErrUnauthorized = echo.NewHTTPError(http.StatusForbidden, "")
	// ErrNotFound : HTTP 404 error
	ErrNotFound = echo.NewHTTPError(http.StatusNotFound, "")
	// ErrBadReqBody : HTTP 400 error
	ErrBadReqBody = echo.NewHTTPError(http.StatusBadRequest, "")
	// ErrGatewayTimeout : HTTP 504 error
	ErrGatewayTimeout = echo.NewHTTPError(http.StatusGatewayTimeout, "")
	// ErrInternal : HTTP 500 error
	ErrInternal = echo.NewHTTPError(http.StatusInternalServerError, "")
	// ErrNotImplemented : HTTP 405 error
	ErrNotImplemented = echo.NewHTTPError(http.StatusNotImplemented, "")
	// ErrExists : HTTP Error
	ErrExists = echo.NewHTTPError(http.StatusSeeOther, "")
)

// Respond : manage responses
func Respond(c echo.Context, st int, b []byte) error {
	if st == 200 {
		return c.JSONBlob(st, b)
	}

	return echo.NewHTTPError(st, string(b))
}

// ErrMessage prepares a message string to be responded
func ErrMessage(msg string) []byte {
	return []byte(`{"message": "` + msg + `"}`)
}
