package gin_unit_test

import "errors"

var (
	// when router is not set
	// 当路由未设置的时候
	ErrRouterNotSet = errors.New("router not set")
	// when the method should be post or put
	// 当请求方式应该为post或put的时候
	ErrMustPostOrPut = errors.New("method must be post or put")
)
