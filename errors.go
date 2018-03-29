package gin_unit_test

import "errors"

var (
	ErrRouterNotSet      = errors.New("router not set")
	ErrMustPostOrPut     = errors.New("method must be post or put")
	ErrMustBeStructOrMap = errors.New("param's type must be struct or map")
)
