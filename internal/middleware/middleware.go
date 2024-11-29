package middleware

import "github.com/zhenyanesterkova/metricsmonitor/internal/app/server/logger"

type MiddlewareStruct struct {
	Logger logger.LogrusLogger
}

func NewMiddlewareStruct(log logger.LogrusLogger) MiddlewareStruct {
	return MiddlewareStruct{
		Logger: log,
	}
}
