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

func IsServiceRunning(err error) bool {
	_, ok := err.(serviceRunningError)
	return ok
}
