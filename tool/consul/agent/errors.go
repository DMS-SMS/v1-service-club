package agent

import "errors"

var (
	ErrAvailableNodeNotFound = errors.New("there is no currently available services")
)
