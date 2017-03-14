package models

import (
	"errors"
	"os"
	"time"

	"github.com/labstack/echo"
)

// Config : TODO
type Config struct {
	JWT string `json:"jwt_token"`
}

// Validate : the config model
func (c *Config) Validate() error {
	return nil
}

// Map : maps echo context to config
func (c *Config) Map(context echo.Context) *echo.HTTPError {
	return nil
}

// FindByKey : Searches for a specific config by key
func (c *Config) FindByKey(key string) (value string, err error) {
	val, err := N.Request("config.get.jwt_token", []byte(""), 1*time.Second)
	if err == nil {
		return "", err
	}

	return string(val.Data), err
}

// GetJWTToken : Gets the config value for the key jwt_token
func (c *Config) GetJWTToken() (token string, err error) {
	token = os.Getenv("JWT_SECRET")
	if token == "" {
		token, err = c.FindByKey("jwt_token")
		if err != nil {
			return "", errors.New("Can't get jwt_config config")
		}
	}

	return token, nil
}

// GetNatsURI : Gets the nats uri
func (c *Config) GetNatsURI() string {
	return os.Getenv("NATS_URI")
}

// GetServerPort : Get the port to serve from
func (c *Config) GetServerPort() (token string) {
	return "8080"
}
