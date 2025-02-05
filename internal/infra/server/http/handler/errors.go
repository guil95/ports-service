package handler

import "errors"

var (
	methodNotAllowed   = errors.New("method not allowed")
	notFound           = errors.New("not found")
	invalidRequest     = errors.New("invalid request")
	internalServer     = errors.New("internal server error")
	missingIDParameter = errors.New("missing id parameter")
)
