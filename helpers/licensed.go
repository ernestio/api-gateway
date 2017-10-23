package helpers

import (
	"errors"
	"os"
)

// Licensed : Checks if the current api is running with premium support
func Licensed() error {
	if os.Getenv("ERNEST_EDITION") != "enterprise" {
		return errors.New("You're running ernest community edition, please contact R3Labs for premium support")
	}
	return nil
}
