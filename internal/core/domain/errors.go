package domain

import "errors"

var ErrPortNotFound = errors.New("port not found")
var ErrInvalidPort = errors.New("invalid port")
var ErrInvalidJson = errors.New("invalid json")
