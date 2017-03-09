package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/labstack/echo"
)

// Usage : Usage-store entity
type Usage struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Service string `json:"service"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	From    int64  `json:"from"`
	To      int64  `json:"to"`
}

// Validate : validates the usage
func (l *Usage) Validate() error {
	if l.Type == "" {
		return errors.New("Usage type is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (l *Usage) Map(c echo.Context) *echo.HTTPError {
	body := c.Request().Body
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return ErrBadReqBody
	}

	err = json.Unmarshal(data, &l)
	if err != nil {
		return ErrBadReqBody
	}

	return nil
}

// FindAll : Searches for all usaages on the system
func (l *Usage) FindAll(usages *[]Usage) (err error) {
	query := make(map[string]interface{})
	if err := NewBaseModel("usage").FindBy(query, usages); err != nil {
		return err
	}
	return nil
}

// FindAllInRange : Searches for all usaages on a date range
func (l *Usage) FindAllInRange(from, to int64, usages *[]Usage) (err error) {
	query := make(map[string]interface{})
	query["from"] = from
	query["to"] = to
	if err := NewBaseModel("usage").FindBy(query, usages); err != nil {
		return err
	}
	return nil
}
