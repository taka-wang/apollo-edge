package dispatcher

import "errors"

var (
	// ErrLoadRouteFile fail to open the routing file
	ErrLoadRouteFile = errors.New("failed to load routing file")

	// ErrRouteFormat the format of route file is incorrect
	ErrRouteFormat = errors.New("the format of routing file is incorrect")
)
