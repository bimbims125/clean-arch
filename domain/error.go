package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInternalServer = errors.New("internal server error")
var ErrBadRequest = errors.New("bad request")
