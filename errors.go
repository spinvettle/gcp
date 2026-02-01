package gcp

import "errors"

var (
	InvalidCapacityError = errors.New("GcpError:Invalid capacity,it must greater than 0")
	PoolFullError        = errors.New("GcpError:Number of workers has retached capacity")
	PoolClosedError      = errors.New("GcpError:Pool already closed")
)
