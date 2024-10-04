package metricerrors

import "errors"

var (
	ErrParseValue    = errors.New("can not parse metric value")
	ErrUnknownType   = errors.New("unknown metric type")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidType   = errors.New("invalid type")
	ErrUnknownMetric = errors.New("unknown metric name")
)
