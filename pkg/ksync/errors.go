package ksync

import (
	"fmt"
)

type serviceRunningError struct {
	error
	service *Service
}

func (e serviceRunningError) Error() string {
	return fmt.Sprintf("Error: Already running: %s", e.service)
}

// IsServiceRunning checks to see if the error is related to a service running.
func IsServiceRunning(err error) bool {
	_, ok := err.(serviceRunningError)
	return ok
}
