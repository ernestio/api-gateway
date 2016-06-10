package main

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	notFound       = echo.NewHTTPError(http.StatusNotFound, "")
	badReqBody     = echo.NewHTTPError(http.StatusBadRequest, "")
	gatewayTimeout = echo.NewHTTPError(http.StatusGatewayTimeout, "")
)
