package controllers

import (
	"github.com/ernestio/api-gateway/controllers/groups"
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/labstack/echo"
)

type list func(models.User) (int, []byte)
type get func(models.User, string) (int, []byte)
type create func(models.User, []byte) (int, []byte)
type update func(models.User, string, []byte) (int, []byte)
type delete func(models.User, string) (int, []byte)

func genericList(c echo.Context, entity string, fn list) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, entity+"s/list")
	if st == 200 {
		st, b = fn(au)
	}

	return h.Respond(c, st, b)
}

func genericGet(c echo.Context, entity string, fn get) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, entity+"s/get")
	if st == 200 {
		g := c.Param(entity)
		st, b = groups.Get(au, g)
	}

	return h.Respond(c, st, b)
}

func genericCreate(c echo.Context, entity string, fn create) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, entity+"s/create")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = fn(au, body)
	}

	return h.Respond(c, st, b)
}

func genericUpdate(c echo.Context, entity string, fn update) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, entity+"s/update")
	if st != 200 {
		return h.Respond(c, st, b)
	}

	st = 500
	b = []byte("Invalid input")
	name := c.Param(entity)

	body, err := h.GetRequestBody(c)
	if err == nil {
		st, b = fn(au, name, body)
	}

	return h.Respond(c, st, b)
}

func genericDelete(c echo.Context, entity string, fn delete) error {
	au := AuthenticatedUser(c)
	st, b := h.IsAuthorized(&au, entity+"s/delete")
	if st == 200 {
		g := c.Param(entity)
		st, b = fn(au, g)
	}

	return h.Respond(c, st, b)
}
