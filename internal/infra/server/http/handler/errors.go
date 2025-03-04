package handler

import "errors"

var (
	invalidRequest     = errors.New("invalid request")
	internalServer     = errors.New("internal server error")
	missingIDParameter = errors.New("missing id parameter")
)
